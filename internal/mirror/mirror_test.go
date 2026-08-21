package mirror

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestSourceImagesFromLock(t *testing.T) {
	lock := Lock{Entries: []LockEntry{
		{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"},
		{SourceRepository: "gvm-tools", SourceChannel: "gvm-tools"},
	}}
	want := []SourceImage{
		{Repository: "gvmd", Channel: "gvmd:stable"},
		{Repository: "gvm-tools", Channel: "gvm-tools"},
	}
	got := SourceImagesFromLock(lock)
	if len(got) != len(want) {
		t.Fatalf("got %d images, want %d: %+v", len(got), len(want), got)
	}
	for i, img := range got {
		if img != want[i] {
			t.Errorf("images[%d] = %+v, want %+v", i, img, want[i])
		}
	}
}

func TestRewriteCompose(t *testing.T) {
	compose := strings.Join([]string{
		"services:",
		"  gvmd:",
		"    image: registry.community.greenbone.net/community/gvmd:stable",
		"  gvm-tools:",
		"    image: ghcr.io/brennoo/gvm-mirror/gvm-tools@sha256:oldcccc # sha256-oldcccc",
		"  unrelated:",
		"    image: postgres:16",
	}, "\n")
	entries := []LockEntry{
		{
			SourceRepository:      "registry.community.greenbone.net/community/gvmd",
			DestinationRepository: "ghcr.io/brennoo/gvm-mirror/gvmd",
			DestinationDigest:     "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:aaaa",
			Version:               "23.0.0",
		},
		{
			SourceRepository:      "gvm-tools",
			DestinationRepository: "ghcr.io/brennoo/gvm-mirror/gvm-tools",
			DestinationDigest:     "ghcr.io/brennoo/gvm-mirror/gvm-tools@sha256:cccc",
			Tags:                  []string{"sha256-cccc"},
		},
	}

	got, err := RewriteCompose([]byte(compose), entries)
	if err != nil {
		t.Fatalf("RewriteCompose: %v", err)
	}
	want := strings.Join([]string{
		"services:",
		"  gvmd:",
		"    image: ghcr.io/brennoo/gvm-mirror/gvmd@sha256:aaaa # 23.0.0",
		"  gvm-tools:",
		"    image: ghcr.io/brennoo/gvm-mirror/gvm-tools@sha256:cccc # sha256-cccc",
		"  unrelated:",
		"    image: postgres:16",
	}, "\n")
	if string(got) != want {
		t.Errorf("RewriteCompose =\n%s\nwant\n%s", got, want)
	}
}

func TestRewriteComposeMissingDestinationDigestFailsClosed(t *testing.T) {
	compose := "    image: gvmd:stable"
	entries := []LockEntry{{SourceRepository: "gvmd", Version: "23.0.0"}}
	if _, err := RewriteCompose([]byte(compose), entries); err == nil {
		t.Fatal("expected an error for a missing destination digest, got nil")
	}
}

// fakeRunner dispatches on the joined command line; missing entries cause a
// test failure so unexpected commands are caught rather than silently
// returning a zero value.
func fakeRunner(t *testing.T, responses map[string]func() ([]byte, error)) Runner {
	t.Helper()
	return func(_ context.Context, name string, args ...string) ([]byte, error) {
		key := name + " " + strings.Join(args, " ")
		resp, ok := responses[key]
		if !ok {
			t.Fatalf("unexpected command: %s", key)
		}
		return resp()
	}
}

func ok(data []byte) func() ([]byte, error) {
	return func() ([]byte, error) { return data, nil }
}

func fail(err error) func() ([]byte, error) {
	return func() ([]byte, error) { return nil, err }
}

func TestResolveDigest(t *testing.T) {
	manifest := []byte(`{"mediaType":"single"}`)
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable": ok(manifest),
	})
	digestRef, err := ResolveDigest(context.Background(), run, "gvmd:stable")
	if err != nil {
		t.Fatalf("ResolveDigest: %v", err)
	}
	sum := sha256.Sum256(manifest)
	want := "gvmd@sha256:" + hex.EncodeToString(sum[:])
	if digestRef != want {
		t.Errorf("ResolveDigest = %q, want %q", digestRef, want)
	}
}

func TestResolveDigestPropagatesError(t *testing.T) {
	wantErr := errors.New("boom")
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable": fail(wantErr),
	})
	if _, err := ResolveDigest(context.Background(), run, "gvmd:stable"); !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want it to wrap %v", err, wantErr)
	}
}

func manifestListJSON(entries ...[3]string) []byte {
	type platform struct {
		OS   string `json:"os"`
		Arch string `json:"architecture"`
	}
	type manifest struct {
		Digest   string   `json:"digest"`
		Platform platform `json:"platform"`
	}
	var list struct {
		Manifests []manifest `json:"manifests"`
	}
	for _, e := range entries {
		list.Manifests = append(list.Manifests, manifest{Digest: e[0], Platform: platform{OS: e[1], Arch: e[2]}})
	}
	data, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}
	return data
}

func TestInspectPlatformsSinglePlatform(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd@sha256:aaaa": ok([]byte(`{"mediaType":"single"}`)),
		"skopeo inspect docker://gvmd@sha256:aaaa":       ok([]byte(`{"Os":"linux","Architecture":"amd64"}`)),
	})
	platforms, err := InspectPlatforms(context.Background(), run, "gvmd@sha256:aaaa")
	if err != nil {
		t.Fatalf("InspectPlatforms: %v", err)
	}
	if len(platforms) != 1 || platforms[0] != (PlatformManifest{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}) {
		t.Errorf("platforms = %+v, want a single linux/amd64 entry for sha256:aaaa", platforms)
	}
}

func TestInspectPlatformsSingleUnknownIdentityFailsClosed(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd@sha256:aaaa": ok([]byte(`{"mediaType":"single"}`)),
		"skopeo inspect docker://gvmd@sha256:aaaa":       ok([]byte(`{}`)),
	})
	if _, err := InspectPlatforms(context.Background(), run, "gvmd@sha256:aaaa"); err == nil {
		t.Fatal("expected an error when os/architecture are empty, got nil")
	}
}

func TestInspectPlatformsManifestList(t *testing.T) {
	data := manifestListJSON(
		[3]string{"sha256:amd64digest", "linux", "amd64"},
		[3]string{"sha256:arm64digest", "linux", "arm64"},
		[3]string{"sha256:attestation", "unknown", "unknown"},
	)
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd@sha256:aaaa": ok(data),
	})
	platforms, err := InspectPlatforms(context.Background(), run, "gvmd@sha256:aaaa")
	if err != nil {
		t.Fatalf("InspectPlatforms: %v", err)
	}
	if len(platforms) != 2 {
		t.Fatalf("platforms = %+v, want 2 real platforms (attestation entry dropped)", platforms)
	}
}

func TestInspectPlatformsManifestListDuplicateFailsClosed(t *testing.T) {
	data := manifestListJSON(
		[3]string{"sha256:amd64digest", "linux", "amd64"},
		[3]string{"sha256:otherdigest", "linux", "amd64"},
	)
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd@sha256:aaaa": ok(data),
	})
	if _, err := InspectPlatforms(context.Background(), run, "gvmd@sha256:aaaa"); err == nil {
		t.Fatal("expected an error for duplicate platform identities, got nil")
	}
}

func TestInspectPlatformsAllUnknownFailsClosed(t *testing.T) {
	data := manifestListJSON([3]string{"sha256:attestation", "unknown", "unknown"})
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd@sha256:aaaa": ok(data),
	})
	if _, err := InspectPlatforms(context.Background(), run, "gvmd@sha256:aaaa"); err == nil {
		t.Fatal("expected an error when no platform entries are real, got nil")
	}
}

func TestPlatformLabels(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect docker://gvmd@sha256:amd64digest": ok([]byte(`{"Labels":{"org.opencontainers.image.version":"23.0.0"}}`)),
	})
	labels, err := PlatformLabels(context.Background(), run, "gvmd", "sha256:amd64digest")
	if err != nil {
		t.Fatalf("PlatformLabels: %v", err)
	}
	if labels["org.opencontainers.image.version"] != "23.0.0" {
		t.Errorf("labels = %+v, missing expected version", labels)
	}
}

func TestResolveVersionRevisionAgreement(t *testing.T) {
	labels := map[string]map[string]string{
		"linux/amd64": {versionLabel: "23.0.0", revisionLabel: "abc123"},
		"linux/arm64": {versionLabel: "23.0.0", revisionLabel: "abc123"},
	}
	version, revision, err := ResolveVersionRevision(labels)
	if err != nil {
		t.Fatalf("ResolveVersionRevision: %v", err)
	}
	if version != "23.0.0" || revision != "abc123" {
		t.Errorf("got (%q, %q), want (23.0.0, abc123)", version, revision)
	}
}

func TestResolveVersionRevisionMissingFailsClosed(t *testing.T) {
	labels := map[string]map[string]string{
		"linux/amd64": {versionLabel: "23.0.0", revisionLabel: "abc123"},
		"linux/arm64": {versionLabel: "23.0.0"},
	}
	if _, _, err := ResolveVersionRevision(labels); err == nil {
		t.Fatal("expected an error for a missing revision label, got nil")
	}
}

func TestResolveVersionRevisionMismatchFailsClosed(t *testing.T) {
	labels := map[string]map[string]string{
		"linux/amd64": {versionLabel: "23.0.0", revisionLabel: "abc123"},
		"linux/arm64": {versionLabel: "23.0.1", revisionLabel: "abc123"},
	}
	if _, _, err := ResolveVersionRevision(labels); err == nil {
		t.Fatal("expected an error for disagreeing version labels, got nil")
	}
}

func TestFallbackTag(t *testing.T) {
	if got, want := FallbackTag("gvmd@sha256:aaaa"), "sha256-aaaa"; got != want {
		t.Errorf("FallbackTag = %q, want %q", got, want)
	}
}

func TestDestinationRepository(t *testing.T) {
	tests := []struct{ source, want string }{
		{"registry.community.greenbone.net/community/gvmd", "ghcr.io/brennoo/gvm-mirror/gvmd"},
		{"gvm-tools", "ghcr.io/brennoo/gvm-mirror/gvm-tools"},
	}
	for _, tt := range tests {
		if got := DestinationRepository("ghcr.io/brennoo/gvm-mirror", tt.source); got != tt.want {
			t.Errorf("DestinationRepository(%q) = %q, want %q", tt.source, got, tt.want)
		}
	}
}

func TestPlanDestinationsCollisionFailsClosed(t *testing.T) {
	images := []SourceImage{
		{Repository: "registry.a.example/community/gvmd", Channel: "x"},
		{Repository: "registry.b.example/other/gvmd", Channel: "y"},
	}
	if _, err := PlanDestinations("ghcr.io/brennoo/gvm-mirror", images); err == nil {
		t.Fatal("expected a collision error, got nil")
	}
}

func TestPlanDestinationsNoCollision(t *testing.T) {
	images := []SourceImage{
		{Repository: "registry.a.example/community/gvmd", Channel: "x"},
		{Repository: "registry.a.example/community/gvm-tools", Channel: "y"},
	}
	dest, err := PlanDestinations("ghcr.io/brennoo/gvm-mirror", images)
	if err != nil {
		t.Fatalf("PlanDestinations: %v", err)
	}
	if len(dest) != 2 {
		t.Fatalf("dest = %+v, want 2 entries", dest)
	}
}

func TestCopy(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:aaaa docker://ghcr.io/x/gvmd:sha256-aaaa": ok(nil),
	})
	if err := Copy(context.Background(), run, "gvmd@sha256:aaaa", "ghcr.io/x/gvmd:sha256-aaaa"); err != nil {
		t.Fatalf("Copy: %v", err)
	}
}

func TestWriteOnceTagCreatesMissingTag(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://ghcr.io/x/gvmd:stable":                                           fail(errors.New("manifest unknown")),
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:aaaa docker://ghcr.io/x/gvmd:stable": ok(nil),
	})
	if err := WriteOnceTag(context.Background(), run, "ghcr.io/x/gvmd", "stable", "gvmd@sha256:aaaa", nil); err != nil {
		t.Fatalf("WriteOnceTag: %v", err)
	}
}

func TestWriteOnceTagNoOpWhenAlreadyCorrect(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://ghcr.io/x/gvmd:stable": ok([]byte(`{"mediaType":"single"}`)),
	})
	sum := sha256.Sum256([]byte(`{"mediaType":"single"}`))
	digestRef := "gvmd@sha256:" + hex.EncodeToString(sum[:])
	if err := WriteOnceTag(context.Background(), run, "ghcr.io/x/gvmd", "stable", digestRef, nil); err != nil {
		t.Fatalf("WriteOnceTag: %v", err)
	}
}

func TestWriteOnceTagConflictFailsClosed(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://ghcr.io/x/gvmd:stable": ok([]byte(`{"mediaType":"different"}`)),
	})
	err := WriteOnceTag(context.Background(), run, "ghcr.io/x/gvmd", "stable", "gvmd@sha256:aaaa", nil)
	if err == nil {
		t.Fatal("expected a write-once conflict error, got nil")
	}
}

func TestWriteOnceTagMovableTagIsRepointed(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://ghcr.io/x/gvmd:latest":                                           ok([]byte(`{"mediaType":"different"}`)),
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:aaaa docker://ghcr.io/x/gvmd:latest": ok(nil),
	})
	err := WriteOnceTag(context.Background(), run, "ghcr.io/x/gvmd", "latest", "gvmd@sha256:aaaa", map[string]bool{"latest": true})
	if err != nil {
		t.Fatalf("WriteOnceTag: %v", err)
	}
}

func TestWriteOnceTagRejectsInvalidTagName(t *testing.T) {
	run := fakeRunner(t, map[string]func() ([]byte, error){})
	if err := WriteOnceTag(context.Background(), run, "ghcr.io/x/gvmd", "-bad", "gvmd@sha256:aaaa", nil); err == nil {
		t.Fatal("expected an error for an invalid tag name, got nil")
	}
}

func TestVerifyPlatformsMatchAgrees(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}, {OS: "linux", Arch: "arm64", Digest: "sha256:bbbb"}}
	got := []PlatformManifest{{OS: "linux", Arch: "arm64", Digest: "sha256:bbbb"}, {OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}}
	if err := VerifyPlatformsMatch(want, got); err != nil {
		t.Fatalf("VerifyPlatformsMatch: %v", err)
	}
}

func TestVerifyPlatformsMatchCountMismatchFailsClosed(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}, {OS: "linux", Arch: "arm64", Digest: "sha256:bbbb"}}
	got := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}}
	if err := VerifyPlatformsMatch(want, got); err == nil {
		t.Fatal("expected an error for a platform count mismatch, got nil")
	}
}

func TestVerifyPlatformsMatchMissingPlatformFailsClosed(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}}
	got := []PlatformManifest{{OS: "linux", Arch: "arm64", Digest: "sha256:aaaa"}}
	if err := VerifyPlatformsMatch(want, got); err == nil {
		t.Fatal("expected an error for a missing platform, got nil")
	}
}

func TestVerifyPlatformsMatchDuplicateWantFailsClosed(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:bbbb"}, {OS: "linux", Arch: "amd64", Digest: "sha256:bbbb"}}
	got := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}, {OS: "linux", Arch: "amd64", Digest: "sha256:bbbb"}}
	if err := VerifyPlatformsMatch(want, got); err == nil {
		t.Fatal("expected an error for a duplicate platform identity in want, got nil")
	}
}

func TestVerifyPlatformsMatchDuplicateGotFailsClosed(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}, {OS: "linux", Arch: "arm64", Digest: "sha256:bbbb"}}
	got := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}, {OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}}
	if err := VerifyPlatformsMatch(want, got); err == nil {
		t.Fatal("expected an error for a duplicate platform identity in got, got nil")
	}
}

func TestVerifyPlatformsMatchDigestMismatchFailsClosed(t *testing.T) {
	want := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:aaaa"}}
	got := []PlatformManifest{{OS: "linux", Arch: "amd64", Digest: "sha256:bbbb"}}
	if err := VerifyPlatformsMatch(want, got); err == nil {
		t.Fatal("expected an error for a digest mismatch, got nil")
	}
}

func TestWriteLockRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/lock.json"
	lock := Lock{Entries: []LockEntry{{
		SourceRepository:      "gvmd",
		SourceChannel:         "gvmd:stable",
		SourceDigest:          "gvmd@sha256:aaaa",
		DestinationRepository: "ghcr.io/x/gvmd",
		DestinationDigest:     "ghcr.io/x/gvmd@sha256:aaaa",
		Version:               "23.0.0",
		Revision:              "abc123",
		Tags:                  []string{"sha256-aaaa", "23.0.0"},
		Platforms:             []string{"linux/amd64"},
	}}}
	if err := WriteLock(path, lock); err != nil {
		t.Fatalf("WriteLock: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading lock file: %v", err)
	}
	var got Lock
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshaling lock file: %v", err)
	}
	if len(got.Entries) != 1 || got.Entries[0].SourceRepository != "gvmd" {
		t.Errorf("got = %+v, want the entry written above", got)
	}
}
