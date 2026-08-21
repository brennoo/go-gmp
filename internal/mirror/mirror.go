// Package mirror copies the images recorded in a lock file into a registry
// this project controls, using Skopeo, records what it did back into that
// lock file, and rewrites a Docker Compose file's pins to match. It exists
// because content-digest pinning alone only protects against a tag moving —
// it does nothing if the upstream registry garbage-collects the blob,
// changes retention, removes the repository, or is simply unreachable.
//
// Every function that shells out takes a Runner so it can be unit-tested
// against a fake instead of a real Skopeo binary and network.
package mirror

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/brennoo/go-gmp/internal/imageref"
)

// Runner executes name with args and returns its stdout, or an error
// embedding stderr — the seam every Skopeo invocation in this package goes
// through.
type Runner func(ctx context.Context, name string, args ...string) ([]byte, error)

// Exec is the real Runner, shelling out via os/exec.
func Exec(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...) //nolint:gosec // name/args are this package's own fixed skopeo invocations
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return []byte(stdout.String()), nil
}

// SourceImage is one distinct image repository to mirror, together with the
// mutable channel (repo:tag, or a bare repo for the implicit-latest case) it
// should be re-resolved from.
type SourceImage struct {
	Repository string
	Channel    string
}

// SourceImagesFromLock returns one SourceImage per entry of a previous run's
// lock file: once hack/docker-compose.yml points at the mirror, its "image:"
// lines no longer carry the upstream channel to re-resolve.
func SourceImagesFromLock(lock Lock) []SourceImage {
	images := make([]SourceImage, len(lock.Entries))
	for i, e := range lock.Entries {
		images[i] = SourceImage{Repository: e.SourceRepository, Channel: e.SourceChannel}
	}
	return images
}

// RewriteCompose rewrites every "    image: ..." line of composeData whose
// repository matches a LockEntry's SourceRepository (cutting over a Compose
// file still on the original upstream) or DestinationRepository (keeping it
// pinned on every run after). A line matching no entry passes through
// unchanged. The replacement tag is the entry's resolved Version if any,
// else its digest-derived fallback tag.
func RewriteCompose(composeData []byte, entries []LockEntry) ([]byte, error) {
	byRepo := make(map[string]LockEntry, len(entries)*2)
	for _, e := range entries {
		byRepo[e.SourceRepository] = e
		byRepo[e.DestinationRepository] = e
	}

	lines := strings.Split(string(composeData), "\n")
	for i, line := range lines {
		repo, ok := imageref.ImageLineRepository(line)
		if !ok {
			continue
		}
		e, ok := byRepo[repo]
		if !ok {
			continue
		}
		if e.DestinationDigest == "" {
			return nil, fmt.Errorf("mirror: line %d: lock entry for %q has no destination digest", i+1, repo)
		}
		tag := e.Version
		if tag == "" {
			if len(e.Tags) == 0 {
				return nil, fmt.Errorf("mirror: line %d: lock entry for %q has no tags", i+1, repo)
			}
			tag = e.Tags[0]
		}
		lines[i] = fmt.Sprintf("    image: %s # %s", e.DestinationDigest, tag)
	}
	return []byte(strings.Join(lines, "\n")), nil
}

func splitDigestRef(digestRef string) (hexDigest string, ok bool) {
	_, rest, found := strings.Cut(digestRef, "@sha256:")
	if !found || rest == "" {
		return "", false
	}
	return rest, true
}

func manifestDigest(ctx context.Context, run Runner, ref string) (string, error) {
	out, err := run(ctx, "skopeo", "inspect", "--raw", "docker://"+ref)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(out)
	return hex.EncodeToString(sum[:]), nil
}

// ResolveDigest resolves ref (a mutable channel or anything Skopeo accepts)
// to its content digest, by hashing the raw manifest bytes rather than
// trusting a registry-reported digest header.
func ResolveDigest(ctx context.Context, run Runner, ref string) (digestRef string, err error) {
	hexDigest, err := manifestDigest(ctx, run, ref)
	if err != nil {
		return "", fmt.Errorf("mirror: resolving digest for %s: %w", ref, err)
	}
	return imageref.Repository(ref) + "@sha256:" + hexDigest, nil
}

// PlatformManifest is one platform's entry in a (possibly single-entry)
// manifest list.
type PlatformManifest struct {
	OS     string
	Arch   string
	Digest string // "sha256:<hex>" per-platform manifest digest
}

// InspectPlatforms returns one PlatformManifest per real platform digestRef's
// manifest advertises. A manifest list/index yields one entry per
// manifests[] item; a single-platform manifest yields one synthetic entry
// carrying digestRef's own digest. Manifest-list entries with no real
// platform (buildkit provenance/attestation entries, which report
// "unknown/unknown") are dropped rather than surfaced as a platform.
func InspectPlatforms(ctx context.Context, run Runner, digestRef string) ([]PlatformManifest, error) {
	out, err := run(ctx, "skopeo", "inspect", "--raw", "docker://"+digestRef)
	if err != nil {
		return nil, fmt.Errorf("mirror: inspecting manifest for %s: %w", digestRef, err)
	}

	var probe struct {
		Manifests []struct {
			Digest   string `json:"digest"`
			Platform struct {
				OS           string `json:"os"`
				Architecture string `json:"architecture"`
			} `json:"platform"`
		} `json:"manifests"`
	}
	if err := json.Unmarshal(out, &probe); err != nil {
		return nil, fmt.Errorf("mirror: parsing manifest for %s: %w", digestRef, err)
	}

	if len(probe.Manifests) == 0 {
		hexDigest, ok := splitDigestRef(digestRef)
		if !ok {
			return nil, fmt.Errorf("mirror: %q is not a resolved image digest", digestRef)
		}
		osName, arch, err := inspectOSArch(ctx, run, digestRef)
		if err != nil {
			return nil, err
		}
		return []PlatformManifest{{OS: osName, Arch: arch, Digest: "sha256:" + hexDigest}}, nil
	}

	platforms := make([]PlatformManifest, 0, len(probe.Manifests))
	for _, m := range probe.Manifests {
		if m.Platform.OS == "unknown" || m.Platform.Architecture == "unknown" {
			continue
		}
		platforms = append(platforms, PlatformManifest{OS: m.Platform.OS, Arch: m.Platform.Architecture, Digest: m.Digest})
	}
	if len(platforms) == 0 {
		return nil, fmt.Errorf("mirror: manifest list for %s has no real platform entries", digestRef)
	}
	if _, err := platformsByKey(platforms); err != nil {
		return nil, fmt.Errorf("mirror: manifest list for %s has %w", digestRef, err)
	}
	return platforms, nil
}

// inspectOSArch returns a single-platform manifest's own OS/architecture, for
// the InspectPlatforms case where there is no manifests[] list to read a
// platform from.
func inspectOSArch(ctx context.Context, run Runner, ref string) (osName, arch string, err error) {
	out, err := run(ctx, "skopeo", "inspect", "docker://"+ref)
	if err != nil {
		return "", "", fmt.Errorf("mirror: inspecting platform for %s: %w", ref, err)
	}
	var parsed struct {
		Os           string `json:"Os"`
		Architecture string `json:"Architecture"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return "", "", fmt.Errorf("mirror: parsing platform for %s: %w", ref, err)
	}
	if parsed.Os == "" || parsed.Architecture == "" {
		return "", "", fmt.Errorf("mirror: %s has no platform identity (os=%q, architecture=%q)", ref, parsed.Os, parsed.Architecture)
	}
	return parsed.Os, parsed.Architecture, nil
}

// platformKey identifies a platform by OS/Arch, ignoring digest.
func platformKey(p PlatformManifest) string { return p.OS + "/" + p.Arch }

// platformsByKey indexes platforms by platformKey, failing closed on any
// duplicate identity — a valid manifest list or platform set never claims
// the same OS/Arch twice, and a map built by silently overwriting duplicates
// would hide exactly the kind of corruption this exists to catch.
func platformsByKey(platforms []PlatformManifest) (map[string]PlatformManifest, error) {
	byKey := make(map[string]PlatformManifest, len(platforms))
	for _, p := range platforms {
		key := platformKey(p)
		if _, dup := byKey[key]; dup {
			return nil, fmt.Errorf("duplicate platform %s", key)
		}
		byKey[key] = p
	}
	return byKey, nil
}

// VerifyPlatformsMatch fails closed unless got carries exactly the same set
// of (OS, Arch, Digest) triples as want, independent of order. Copy's
// --preserve-digests promises this but does not itself verify it — callers
// must re-inspect the destination and compare.
func VerifyPlatformsMatch(want, got []PlatformManifest) error {
	if len(want) != len(got) {
		return fmt.Errorf("mirror: destination has %d platform(s), source has %d", len(got), len(want))
	}
	if _, err := platformsByKey(want); err != nil {
		return fmt.Errorf("mirror: source %w", err)
	}
	byPlatform, err := platformsByKey(got)
	if err != nil {
		return fmt.Errorf("mirror: destination %w", err)
	}
	for _, w := range want {
		name := platformKey(w)
		g, ok := byPlatform[name]
		if !ok {
			return fmt.Errorf("mirror: destination is missing platform %s", name)
		}
		if g.Digest != w.Digest {
			return fmt.Errorf("mirror: platform %s digest mismatch: source %s, destination %s", name, w.Digest, g.Digest)
		}
	}
	return nil
}

// PlatformLabels returns the image config labels for one platform manifest
// of repo.
func PlatformLabels(ctx context.Context, run Runner, repo, platformDigest string) (map[string]string, error) {
	ref := repo + "@" + platformDigest
	out, err := run(ctx, "skopeo", "inspect", "docker://"+ref)
	if err != nil {
		return nil, fmt.Errorf("mirror: inspecting labels for %s: %w", ref, err)
	}
	var parsed struct {
		Labels map[string]string `json:"Labels"`
	}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return nil, fmt.Errorf("mirror: parsing labels for %s: %w", ref, err)
	}
	return parsed.Labels, nil
}

const (
	versionLabel  = "org.opencontainers.image.version"
	revisionLabel = "org.opencontainers.image.revision"
)

// ResolveVersionRevision requires org.opencontainers.image.version and
// .revision to be present, non-empty, and identical across every platform in
// labelsByPlatform (keyed by an arbitrary but stable per-platform label,
// e.g. "linux/amd64"), and returns their agreed value. It fails closed,
// naming exactly which platform(s) are missing or disagree, rather than
// picking one platform's value.
func ResolveVersionRevision(labelsByPlatform map[string]map[string]string) (version, revision string, err error) {
	if len(labelsByPlatform) == 0 {
		return "", "", errors.New("mirror: no platforms to resolve labels from")
	}
	platforms := make([]string, 0, len(labelsByPlatform))
	for p := range labelsByPlatform {
		platforms = append(platforms, p)
	}
	sort.Strings(platforms)

	version, err = agreeingLabel(labelsByPlatform, platforms, versionLabel)
	if err != nil {
		return "", "", err
	}
	revision, err = agreeingLabel(labelsByPlatform, platforms, revisionLabel)
	if err != nil {
		return "", "", err
	}
	return version, revision, nil
}

func agreeingLabel(labelsByPlatform map[string]map[string]string, platforms []string, label string) (string, error) {
	var value string
	var have bool
	var missing, mismatched []string

	for _, p := range platforms {
		v, ok := labelsByPlatform[p][label]
		if !ok || v == "" {
			missing = append(missing, p)
			continue
		}
		if !have {
			value, have = v, true
			continue
		}
		if v != value {
			mismatched = append(mismatched, p)
		}
	}
	if len(missing) > 0 {
		return "", fmt.Errorf("mirror: label %q is missing on platform(s) %s", label, strings.Join(missing, ", "))
	}
	if len(mismatched) > 0 {
		return "", fmt.Errorf("mirror: label %q disagrees across platforms: %q, but %s differ", label, value, strings.Join(mismatched, ", "))
	}
	return value, nil
}

// FallbackTag derives the always-available destination tag for digestRef,
// which never depends on labels being present or consistent.
func FallbackTag(digestRef string) string {
	hexDigest, _ := splitDigestRef(digestRef)
	return "sha256-" + hexDigest
}

// DestinationRepository maps a source repository to its mirror destination
// under prefix, using the source repository's last path segment as the
// mirror's image name.
func DestinationRepository(prefix, sourceRepository string) string {
	name := sourceRepository
	if i := strings.LastIndexByte(sourceRepository, '/'); i >= 0 {
		name = sourceRepository[i+1:]
	}
	return strings.TrimRight(prefix, "/") + "/" + name
}

// PlanDestinations computes DestinationRepository for every image and fails
// closed if two distinct source repositories would collide on the same
// destination name.
func PlanDestinations(prefix string, images []SourceImage) (map[string]string, error) {
	dest := make(map[string]string, len(images))
	sourceOf := make(map[string]string, len(images))
	for _, img := range images {
		d := DestinationRepository(prefix, img.Repository)
		if existing, ok := sourceOf[d]; ok && existing != img.Repository {
			return nil, fmt.Errorf("mirror: destination repository %q would be shared by %q and %q", d, existing, img.Repository)
		}
		sourceOf[d] = img.Repository
		dest[img.Repository] = d
	}
	return dest, nil
}

// Copy copies srcDigestRef to dstRef, preserving every platform and every
// manifest's exact bytes (--all --preserve-digests) so a media-type
// conversion that would silently change content-addressing fails loudly
// instead of producing a mismatched digest discovered later.
func Copy(ctx context.Context, run Runner, srcDigestRef, dstRef string) error {
	if _, err := run(ctx, "skopeo", "copy", "--all", "--preserve-digests", "docker://"+srcDigestRef, "docker://"+dstRef); err != nil {
		return fmt.Errorf("mirror: copying %s to %s: %w", srcDigestRef, dstRef, err)
	}
	return nil
}

// missingManifestMarkers are substrings Skopeo/registries are known to
// include in an error when a tag or manifest simply does not exist yet.
// Any other error is treated as a hard failure rather than "assume missing"
// — an ambiguous error must not be silently read as permission to create a
// possibly-conflicting tag.
var missingManifestMarkers = []string{
	"manifest unknown",
	"MANIFEST_UNKNOWN",
	"name unknown",
	"NAME_UNKNOWN",
	"not found",
}

func isMissingManifest(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, marker := range missingManifestMarkers {
		if strings.Contains(msg, strings.ToLower(marker)) {
			return true
		}
	}
	return false
}

// WriteOnceTag ensures dstRepo:tag points at wantDigestRef. A missing tag is
// created. A tag already pointing at wantDigestRef's digest is a no-op. A
// tag pointing at a different digest is an error unless tag is in movable,
// in which case it is repointed — write-once destination tags are the
// default, movable tags are an explicit opt-in per tag.
func WriteOnceTag(ctx context.Context, run Runner, dstRepo, tag, wantDigestRef string, movable map[string]bool) error {
	if err := imageref.ValidateTagName(tag); err != nil {
		return fmt.Errorf("mirror: %w", err)
	}
	wantHex, ok := splitDigestRef(wantDigestRef)
	if !ok {
		return fmt.Errorf("mirror: %q is not a resolved image digest", wantDigestRef)
	}
	dstRef := dstRepo + ":" + tag

	currentHex, err := manifestDigest(ctx, run, dstRef)
	switch {
	case err != nil && isMissingManifest(err):
		return Copy(ctx, run, wantDigestRef, dstRef)
	case err != nil:
		return fmt.Errorf("mirror: checking existing tag %s: %w", dstRef, err)
	case currentHex == wantHex:
		return nil
	case movable[tag]:
		return Copy(ctx, run, wantDigestRef, dstRef)
	default:
		return fmt.Errorf("mirror: tag %s already points at sha256:%s, refusing to move it to sha256:%s (not configured as movable)", dstRef, currentHex, wantHex)
	}
}

// LockEntry records everything one mirrored image did in a run.
type LockEntry struct {
	SourceRepository      string    `json:"source_repository"`
	SourceChannel         string    `json:"source_channel"`
	SourceDigest          string    `json:"source_digest"`
	DestinationRepository string    `json:"destination_repository"`
	DestinationDigest     string    `json:"destination_digest"`
	Version               string    `json:"version,omitempty"`
	Revision              string    `json:"revision,omitempty"`
	Tags                  []string  `json:"tags"`
	Platforms             []string  `json:"platforms"`
	CapturedAt            time.Time `json:"captured_at"`
}

// Lock is the top-level shape of the machine-readable mirror lock file.
type Lock struct {
	Entries []LockEntry `json:"entries"`
}

// WriteLock writes lock to path as indented JSON with a trailing newline.
func WriteLock(path string, lock Lock) error {
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644) //nolint:gosec // a lock file is not sensitive
}

// ReadLock reads and parses the lock file at path — the durable seed of what
// a mirror run should re-resolve from.
func ReadLock(path string) (Lock, error) {
	data, err := os.ReadFile(path) //nolint:gosec // path is an operator-supplied CLI flag, not untrusted input
	if err != nil {
		return Lock{}, err
	}
	var lock Lock
	if err := json.Unmarshal(data, &lock); err != nil {
		return Lock{}, fmt.Errorf("parsing lock file %s: %w", path, err)
	}
	return lock, nil
}
