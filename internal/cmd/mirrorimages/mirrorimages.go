package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/brennoo/go-gmp/internal/imageref"
	"github.com/brennoo/go-gmp/internal/mirror"
)

// run parses flags, mirrors every image the lock file's previous run
// recorded, and (outside -dry-run) writes the updated lock file and rewrites
// compose-file's pins to match — even if some images failed, so one broken
// source never costs the rest of the run its progress.
func run(ctx context.Context, args []string, runner mirror.Runner, now func() time.Time, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("mirrorimages", flag.ContinueOnError)
	composeFile := fs.String("compose-file", "hack/docker-compose.yml",
		"the Compose file whose image: pins are rewritten to the mirrored digests")
	lockFile := fs.String("lock-file", "hack/image-mirror.lock.json",
		"the mirror's provenance lock file: read as this run's seed of what to mirror, then overwritten with the result")
	destPrefix := fs.String("dest-prefix", "ghcr.io/brennoo/gvm-mirror",
		"destination repository prefix; each source image is mirrored to <dest-prefix>/<name>")
	dryRun := fs.Bool("dry-run", false,
		"resolve and plan, but skip copying and writing the lock file")
	allowMove := fs.String("allow-move", "",
		"comma-separated destination tag names allowed to move to a new digest (default: none, all tags are write-once)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("unexpected argument %q; mirrorimages takes flags only", fs.Arg(0))
	}
	return runWithFlags(ctx, *composeFile, *lockFile, *destPrefix, *dryRun, parseMovable(*allowMove), runner, now, stdout, stderr)
}

func parseMovable(allowMove string) map[string]bool {
	movable := map[string]bool{}
	for _, tag := range strings.Split(allowMove, ",") {
		if tag = strings.TrimSpace(tag); tag != "" {
			movable[tag] = true
		}
	}
	return movable
}

func runWithFlags(ctx context.Context, composeFile, lockFile, destPrefix string, dryRun bool, movable map[string]bool, runner mirror.Runner, now func() time.Time, stdout, stderr io.Writer) error {
	if !dryRun {
		if err := checkSkopeoAvailable(); err != nil {
			return err
		}
	}

	seedLock, err := mirror.ReadLock(lockFile)
	if err != nil {
		return fmt.Errorf("reading seed lock file %s (mirrorimages needs a prior run's lock file to know what to mirror): %w", lockFile, err)
	}
	images := mirror.SourceImagesFromLock(seedLock)

	destinations, err := mirror.PlanDestinations(destPrefix, images)
	if err != nil {
		return err
	}

	// A source that fails this run keeps its seed entry rather than vanishing
	// from the lock file, so the next run still has something to retry it from.
	byRepo := make(map[string]mirror.LockEntry, len(seedLock.Entries))
	for _, e := range seedLock.Entries {
		byRepo[e.SourceRepository] = e
	}

	var succeeded []mirror.LockEntry
	var failed []string
	for _, img := range images {
		entry, err := mirrorOne(ctx, runner, img, destinations[img.Repository], movable, dryRun, now, stderr)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "mirrorimages: %s: %v\n", img.Repository, err)
			failed = append(failed, img.Repository)
			continue
		}
		byRepo[img.Repository] = entry
		succeeded = append(succeeded, entry)
	}
	sort.Slice(succeeded, func(i, j int) bool { return succeeded[i].SourceRepository < succeeded[j].SourceRepository })

	if dryRun {
		for _, e := range succeeded {
			_, _ = fmt.Fprintf(stdout, "%s -> %s tags=%v platforms=%v\n", e.SourceRepository, e.DestinationRepository, e.Tags, e.Platforms)
		}
	} else if err := writeResults(lockFile, composeFile, byRepo, succeeded, failed, stderr); err != nil {
		return err
	}

	if len(failed) > 0 {
		return fmt.Errorf("failed to mirror %d image(s): %s", len(failed), strings.Join(failed, ", "))
	}
	return nil
}

// writeResults writes byRepo (the seed lock merged with this run's
// successes) to lockFile, then rewrites composeFile's pins for every entry
// that has ever mirrored successfully.
func writeResults(lockFile, composeFile string, byRepo map[string]mirror.LockEntry, succeeded []mirror.LockEntry, failed []string, stderr io.Writer) error {
	entries := make([]mirror.LockEntry, 0, len(byRepo))
	for _, e := range byRepo {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].SourceRepository < entries[j].SourceRepository })

	if err := mirror.WriteLock(lockFile, mirror.Lock{Entries: entries}); err != nil {
		return fmt.Errorf("writing lock file: %w", err)
	}
	_, _ = fmt.Fprintf(stderr, "wrote %s (%d image(s) mirrored, %d failed)\n", lockFile, len(succeeded), len(failed))

	composeData, err := os.ReadFile(composeFile) //nolint:gosec // composeFile is an operator-supplied CLI flag, not untrusted input
	if err != nil {
		return fmt.Errorf("reading compose file to rewrite its pins: %w", err)
	}
	// A source with no DestinationDigest has never mirrored successfully, so
	// it's excluded here rather than handed to RewriteCompose, which would
	// otherwise fail closed on a digest this run was never going to have.
	rewritable := make([]mirror.LockEntry, 0, len(entries))
	for _, e := range entries {
		if e.DestinationDigest != "" {
			rewritable = append(rewritable, e)
		}
	}
	rewritten, err := mirror.RewriteCompose(composeData, rewritable)
	if err != nil {
		return fmt.Errorf("rewriting compose file pins: %w", err)
	}
	if err := os.WriteFile(composeFile, rewritten, 0o644); err != nil { //nolint:gosec // matching the compose file's own pre-existing permissions
		return fmt.Errorf("writing rewritten compose file: %w", err)
	}
	_, _ = fmt.Fprintf(stderr, "rewrote %s\n", composeFile)
	return nil
}

// mirrorOne resolves, inspects, and (outside dryRun) copies and tags one
// source image, returning the LockEntry describing what was done or found.
// A source image whose org.opencontainers.image.version/.revision labels
// don't resolve is not a failure of the whole entry — it mirrors under its
// digest-derived fallback tag only, noted on stderr.
func mirrorOne(ctx context.Context, runner mirror.Runner, img mirror.SourceImage, destRepo string, movable map[string]bool, dryRun bool, now func() time.Time, stderr io.Writer) (mirror.LockEntry, error) {
	srcDigestRef, err := mirror.ResolveDigest(ctx, runner, img.Channel)
	if err != nil {
		return mirror.LockEntry{}, err
	}

	platforms, err := mirror.InspectPlatforms(ctx, runner, srcDigestRef)
	if err != nil {
		return mirror.LockEntry{}, err
	}

	labelsByPlatform := make(map[string]map[string]string, len(platforms))
	platformNames := make([]string, 0, len(platforms))
	for _, p := range platforms {
		labels, err := mirror.PlatformLabels(ctx, runner, img.Repository, p.Digest)
		if err != nil {
			return mirror.LockEntry{}, err
		}
		name := p.OS + "/" + p.Arch
		labelsByPlatform[name] = labels
		platformNames = append(platformNames, name)
	}
	sort.Strings(platformNames)

	fallback := mirror.FallbackTag(srcDigestRef)
	tags := []string{fallback}
	version, revision, verErr := mirror.ResolveVersionRevision(labelsByPlatform)
	switch {
	case verErr != nil:
		_, _ = fmt.Fprintf(stderr, "mirrorimages: %s: no version/revision tags (%v); mirroring under %s only\n", img.Repository, verErr, fallback)
	default:
		if err := imageref.ValidateTagName(version); err != nil {
			return mirror.LockEntry{}, fmt.Errorf("version label %q: %w", version, err)
		}
		if err := imageref.ValidateTagName(revision); err != nil {
			return mirror.LockEntry{}, fmt.Errorf("revision label %q: %w", revision, err)
		}
		tags = append(tags, version, revision)
	}

	entry := mirror.LockEntry{
		SourceRepository:      img.Repository,
		SourceChannel:         img.Channel,
		SourceDigest:          srcDigestRef,
		DestinationRepository: destRepo,
		Version:               version,
		Revision:              revision,
		Tags:                  tags,
		Platforms:             platformNames,
		CapturedAt:            now(),
	}
	if dryRun {
		return entry, nil
	}

	if err := mirror.WriteOnceTag(ctx, runner, destRepo, fallback, srcDigestRef, nil); err != nil {
		return mirror.LockEntry{}, err
	}
	dstDigestRef, err := mirror.ResolveDigest(ctx, runner, destRepo+":"+fallback)
	if err != nil {
		return mirror.LockEntry{}, err
	}
	entry.DestinationDigest = dstDigestRef

	dstPlatforms, err := mirror.InspectPlatforms(ctx, runner, dstDigestRef)
	if err != nil {
		return mirror.LockEntry{}, fmt.Errorf("verifying destination platforms for %s: %w", destRepo, err)
	}
	if err := mirror.VerifyPlatformsMatch(platforms, dstPlatforms); err != nil {
		return mirror.LockEntry{}, err
	}

	for _, tag := range tags[1:] {
		if err := mirror.WriteOnceTag(ctx, runner, destRepo, tag, dstDigestRef, movable); err != nil {
			return mirror.LockEntry{}, err
		}
	}
	return entry, nil
}
