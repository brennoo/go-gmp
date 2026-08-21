package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/brennoo/go-gmp/internal/mirror"
)

// fakeRunner dispatches on the joined command line; an unmatched command
// fails the test rather than silently returning a zero value, so a test
// only "passes" if the run took exactly the commands it expects.
func fakeRunner(t *testing.T, responses map[string]func() ([]byte, error)) mirror.Runner {
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

func ok(data []byte) func() ([]byte, error) { return func() ([]byte, error) { return data, nil } }
func fail(err error) func() ([]byte, error) { return func() ([]byte, error) { return nil, err } }

// sequence returns fns[0] on the first call, fns[1] on the second, and so on
// — for a command whose result legitimately differs across repeated calls,
// such as checking a tag before and after creating it.
func sequence(t *testing.T, fns ...func() ([]byte, error)) func() ([]byte, error) {
	t.Helper()
	i := 0
	return func() ([]byte, error) {
		if i >= len(fns) {
			t.Fatalf("sequence: called %d times, only %d responses configured", i+1, len(fns))
		}
		f := fns[i]
		i++
		return f()
	}
}

func digestHex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func fixedNow() time.Time { return time.Date(2026, 8, 17, 0, 0, 0, 0, time.UTC) }

const gvmdManifest = `{"mediaType":"single"}`
const gvmdInspect = `{"Os":"linux","Architecture":"amd64","Labels":{"org.opencontainers.image.version":"23.0.0","org.opencontainers.image.revision":"abc123"}}`
const gvmdInspectNoLabels = `{"Os":"linux","Architecture":"amd64","Labels":{}}`

func writeCompose(t *testing.T, dir, contents string) string {
	t.Helper()
	path := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("writing compose file: %v", err)
	}
	return path
}

// writeSeedLock writes entries as the lock file mirrorimages reads at
// startup to learn what to mirror from.
func writeSeedLock(t *testing.T, dir string, entries ...mirror.LockEntry) string {
	t.Helper()
	path := filepath.Join(dir, "lock.json")
	if err := mirror.WriteLock(path, mirror.Lock{Entries: entries}); err != nil {
		t.Fatalf("writing seed lock file: %v", err)
	}
	return path
}

func readLock(t *testing.T, path string) mirror.Lock {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading lock file: %v", err)
	}
	var lock mirror.Lock
	if err := json.Unmarshal(data, &lock); err != nil {
		t.Fatalf("unmarshaling lock file: %v", err)
	}
	return lock
}

func TestRunDryRunDoesNotCopyOrWriteLock(t *testing.T) {
	dir := t.TempDir()
	compose := filepath.Join(dir, "docker-compose.yml")
	lockPath := writeSeedLock(t, dir, mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"})

	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":                                    ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + digestHex([]byte(gvmdManifest)): ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + digestHex([]byte(gvmdManifest)):       ok([]byte(gvmdInspect)),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
		"-dry-run",
	}, runner, fixedNow, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "gvmd -> ghcr.io/brennoo/gvm-mirror/gvmd") {
		t.Errorf("stdout = %q, want a plan line for gvmd", stdout.String())
	}
	seed := readLock(t, lockPath)
	if len(seed.Entries) != 1 || seed.Entries[0].DestinationDigest != "" {
		t.Errorf("dry-run must not overwrite the seed lock file, got %+v", seed.Entries)
	}
}

func TestRunMirrorsAndWritesLock(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, "    image: gvmd:stable\n")
	lockPath := writeSeedLock(t, dir, mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"})

	srcHex := digestHex([]byte(gvmdManifest))
	fallbackTag := "sha256-" + srcHex
	dstFallbackRef := "ghcr.io/brennoo/gvm-mirror/gvmd:" + fallbackTag
	dstDigestRef := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + srcHex

	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":                                                                          ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + srcHex:                                                                ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + srcHex:                                                                      ok([]byte(gvmdInspect)),
		"skopeo inspect --raw docker://" + dstFallbackRef:                                                                    sequence(t, fail(errors.New("manifest unknown")), ok([]byte(gvmdManifest))),
		"skopeo inspect --raw docker://" + dstDigestRef:                                                                      ok([]byte(gvmdManifest)),
		"skopeo inspect docker://" + dstDigestRef:                                                                            ok([]byte(gvmdInspect)),
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:" + srcHex + " docker://" + dstFallbackRef:                ok(nil),
		"skopeo inspect --raw docker://ghcr.io/brennoo/gvm-mirror/gvmd:23.0.0":                                               fail(errors.New("manifest unknown")),
		"skopeo copy --all --preserve-digests docker://" + dstDigestRef + " docker://ghcr.io/brennoo/gvm-mirror/gvmd:23.0.0": ok(nil),
		"skopeo inspect --raw docker://ghcr.io/brennoo/gvm-mirror/gvmd:abc123":                                               fail(errors.New("manifest unknown")),
		"skopeo copy --all --preserve-digests docker://" + dstDigestRef + " docker://ghcr.io/brennoo/gvm-mirror/gvmd:abc123": ok(nil),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}

	lock := readLock(t, lockPath)
	if len(lock.Entries) != 1 {
		t.Fatalf("lock.Entries = %+v, want 1 entry", lock.Entries)
	}
	entry := lock.Entries[0]
	if entry.SourceRepository != "gvmd" || entry.Version != "23.0.0" || entry.Revision != "abc123" {
		t.Errorf("entry = %+v, missing expected fields", entry)
	}
	wantTags := []string{fallbackTag, "23.0.0", "abc123"}
	if strings.Join(entry.Tags, ",") != strings.Join(wantTags, ",") {
		t.Errorf("entry.Tags = %v, want %v", entry.Tags, wantTags)
	}
	if entry.DestinationDigest != dstDigestRef {
		t.Errorf("entry.DestinationDigest = %q, want %q", entry.DestinationDigest, dstDigestRef)
	}
	if strings.Join(entry.Platforms, ",") != "linux/amd64" {
		t.Errorf("entry.Platforms = %v, want [linux/amd64]", entry.Platforms)
	}

	composeData, err := os.ReadFile(compose)
	if err != nil {
		t.Fatalf("reading rewritten compose file: %v", err)
	}
	wantCompose := "    image: " + dstDigestRef + " # 23.0.0\n"
	if string(composeData) != wantCompose {
		t.Errorf("compose file = %q, want %q", composeData, wantCompose)
	}
}

func TestRunFallbackTagAlreadyCorrectIsNoOp(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, "    image: gvmd:stable\n")
	lockPath := writeSeedLock(t, dir, mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"})

	srcHex := digestHex([]byte(gvmdManifest))
	fallbackTag := "sha256-" + srcHex
	dstFallbackRef := "ghcr.io/brennoo/gvm-mirror/gvmd:" + fallbackTag
	dstDigestRef := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + srcHex

	// No "skopeo copy" entry for the fallback tag: fakeRunner fails the test
	// if WriteOnceTag copies instead of no-opping on an already-correct tag.
	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":           ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + srcHex: ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + srcHex:       ok([]byte(gvmdInspectNoLabels)),
		"skopeo inspect --raw docker://" + dstFallbackRef:     ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://" + dstDigestRef:       ok([]byte(gvmdManifest)),
		"skopeo inspect docker://" + dstDigestRef:             ok([]byte(gvmdInspectNoLabels)),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}
}

func TestRunFallbackTagConflictFailsClosed(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, "    image: gvmd:stable\n")
	lockPath := writeSeedLock(t, dir, mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"})

	srcHex := digestHex([]byte(gvmdManifest))
	fallbackTag := "sha256-" + srcHex
	dstFallbackRef := "ghcr.io/brennoo/gvm-mirror/gvmd:" + fallbackTag

	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":           ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + srcHex: ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + srcHex:       ok([]byte(gvmdInspectNoLabels)),
		// The fallback tag already exists but points at unrelated content —
		// WriteOnceTag must refuse to move it, not silently overwrite it.
		"skopeo inspect --raw docker://" + dstFallbackRef: ok([]byte(`{"mediaType":"conflicting"}`)),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error when the fallback tag conflicts, got nil")
	}
	if !strings.Contains(stderr.String(), "gvmd") {
		t.Errorf("stderr = %q, want it to mention the failed image", stderr.String())
	}

	lock := readLock(t, lockPath)
	if len(lock.Entries) != 1 || lock.Entries[0].SourceRepository != "gvmd" || lock.Entries[0].DestinationDigest != "" {
		t.Errorf("lock.Entries = %+v, want the seed entry preserved unchanged so the next run can retry it", lock.Entries)
	}
}

func TestRunOneFailurePreservesOthersAndExitsNonZero(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, strings.Join([]string{
		"    image: gvmd:stable",
		"    image: broken:latest",
	}, "\n"))
	lockPath := writeSeedLock(t, dir,
		mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"},
		mirror.LockEntry{SourceRepository: "broken", SourceChannel: "broken:latest"},
	)

	srcHex := digestHex([]byte(gvmdManifest))
	fallbackTag := "sha256-" + srcHex
	dstFallbackRef := "ghcr.io/brennoo/gvm-mirror/gvmd:" + fallbackTag
	dstDigestRef := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + srcHex

	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":                                                                          ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + srcHex:                                                                ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + srcHex:                                                                      ok([]byte(gvmdInspect)),
		"skopeo inspect --raw docker://" + dstFallbackRef:                                                                    sequence(t, fail(errors.New("manifest unknown")), ok([]byte(gvmdManifest))),
		"skopeo inspect --raw docker://" + dstDigestRef:                                                                      ok([]byte(gvmdManifest)),
		"skopeo inspect docker://" + dstDigestRef:                                                                            ok([]byte(gvmdInspect)),
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:" + srcHex + " docker://" + dstFallbackRef:                ok(nil),
		"skopeo inspect --raw docker://ghcr.io/brennoo/gvm-mirror/gvmd:23.0.0":                                               fail(errors.New("manifest unknown")),
		"skopeo copy --all --preserve-digests docker://" + dstDigestRef + " docker://ghcr.io/brennoo/gvm-mirror/gvmd:23.0.0": ok(nil),
		"skopeo inspect --raw docker://ghcr.io/brennoo/gvm-mirror/gvmd:abc123":                                               fail(errors.New("manifest unknown")),
		"skopeo copy --all --preserve-digests docker://" + dstDigestRef + " docker://ghcr.io/brennoo/gvm-mirror/gvmd:abc123": ok(nil),
		"skopeo inspect --raw docker://broken:latest":                                                                        fail(errors.New("boom: connection refused")),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected run to return an error when one image fails, got nil")
	}
	if !strings.Contains(stderr.String(), "broken") {
		t.Errorf("stderr = %q, want it to mention the failed image", stderr.String())
	}

	lock := readLock(t, lockPath)
	if len(lock.Entries) != 2 {
		t.Fatalf("lock.Entries = %+v, want both the preserved broken entry and the mirrored gvmd entry", lock.Entries)
	}
	byRepo := make(map[string]mirror.LockEntry, len(lock.Entries))
	for _, e := range lock.Entries {
		byRepo[e.SourceRepository] = e
	}
	if e := byRepo["broken"]; e.DestinationDigest != "" || e.SourceChannel != "broken:latest" {
		t.Errorf("broken entry = %+v, want its seed entry preserved unchanged so the next run can retry it", e)
	}
	if e := byRepo["gvmd"]; e.DestinationDigest == "" {
		t.Errorf("gvmd entry = %+v, want it mirrored successfully", e)
	}

	composeData, err := os.ReadFile(compose)
	if err != nil {
		t.Fatalf("reading compose file: %v", err)
	}
	if !strings.Contains(string(composeData), "image: "+dstDigestRef) {
		t.Errorf("compose = %q, want gvmd's line rewritten to its mirrored digest", composeData)
	}
	if !strings.Contains(string(composeData), "image: broken:latest") {
		t.Errorf("compose = %q, want broken's line left pointing at its unmirrored upstream channel", composeData)
	}
}

func TestRunDestinationCollisionFailsBeforeAnyCommand(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, strings.Join([]string{
		"    image: registry.a.example/community/gvmd:stable",
		"    image: registry.b.example/other/gvmd:stable",
	}, "\n"))
	lockPath := writeSeedLock(t, dir,
		mirror.LockEntry{SourceRepository: "registry.a.example/community/gvmd", SourceChannel: "registry.a.example/community/gvmd:stable"},
		mirror.LockEntry{SourceRepository: "registry.b.example/other/gvmd", SourceChannel: "registry.b.example/other/gvmd:stable"},
	)
	seedBefore := readLock(t, lockPath)

	runner := fakeRunner(t, map[string]func() ([]byte, error){})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected a destination collision error, got nil")
	}
	seedAfter := readLock(t, lockPath)
	if len(seedAfter.Entries) != len(seedBefore.Entries) || seedAfter.Entries[0].DestinationDigest != "" {
		t.Errorf("a collision must fail before the lock file is overwritten, got %+v", seedAfter.Entries)
	}
}

func TestRunMissingVersionLabelFallsBackToDigestTagOnly(t *testing.T) {
	dir := t.TempDir()
	compose := writeCompose(t, dir, "    image: gvmd:stable\n")
	lockPath := writeSeedLock(t, dir, mirror.LockEntry{SourceRepository: "gvmd", SourceChannel: "gvmd:stable"})

	srcHex := digestHex([]byte(gvmdManifest))
	fallbackTag := "sha256-" + srcHex
	dstFallbackRef := "ghcr.io/brennoo/gvm-mirror/gvmd:" + fallbackTag
	dstDigestRef := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + srcHex

	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://gvmd:stable":                                                           ok([]byte(gvmdManifest)),
		"skopeo inspect --raw docker://gvmd@sha256:" + srcHex:                                                 ok([]byte(gvmdManifest)),
		"skopeo inspect docker://gvmd@sha256:" + srcHex:                                                       ok([]byte(gvmdInspectNoLabels)),
		"skopeo inspect --raw docker://" + dstFallbackRef:                                                     sequence(t, fail(errors.New("manifest unknown")), ok([]byte(gvmdManifest))),
		"skopeo inspect --raw docker://" + dstDigestRef:                                                       ok([]byte(gvmdManifest)),
		"skopeo inspect docker://" + dstDigestRef:                                                             ok([]byte(gvmdInspectNoLabels)),
		"skopeo copy --all --preserve-digests docker://gvmd@sha256:" + srcHex + " docker://" + dstFallbackRef: ok(nil),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(), []string{
		"-compose-file", compose,
		"-lock-file", lockPath,
	}, runner, fixedNow, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v (stderr: %s)", err, stderr.String())
	}
	if !strings.Contains(stderr.String(), "no version/revision tags") {
		t.Errorf("stderr = %q, want a notice about missing version/revision labels", stderr.String())
	}

	lock := readLock(t, lockPath)
	if len(lock.Entries) != 1 {
		t.Fatalf("lock.Entries = %+v, want 1 entry", lock.Entries)
	}
	if got := lock.Entries[0].Tags; len(got) != 1 || got[0] != fallbackTag {
		t.Errorf("Tags = %v, want only the fallback tag %q", got, fallbackTag)
	}
	if lock.Entries[0].Version != "" || lock.Entries[0].Revision != "" {
		t.Errorf("Version/Revision = %q/%q, want both empty", lock.Entries[0].Version, lock.Entries[0].Revision)
	}
}

func TestRunUnexpectedArgument(t *testing.T) {
	err := run(context.Background(), []string{"bogus"}, fakeRunner(t, nil), fixedNow, &bytes.Buffer{}, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected an error for an unexpected positional argument, got nil")
	}
}
