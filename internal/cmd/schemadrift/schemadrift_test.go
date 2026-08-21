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

	gmp "github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/internal/codegen"
	"github.com/brennoo/go-gmp/internal/schema"
)

const testSourceRepository = "registry.community.greenbone.net/community/gvmd"

// fakeRunner dispatches on the joined command line; an unlisted command
// fails the test rather than silently returning a zero value.
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

func respond(data []byte) func() ([]byte, error) {
	return func() ([]byte, error) { return data, nil }
}

// writeLockFile writes a fixture lock file directly, rather than through the
// mirror's own WriteLock: schemadrift only ever reads this file, and the
// real writer lives in github.com/brennoo/gvm-mirror now.
func writeLockFile(t *testing.T, dir string, entries ...lockEntry) string {
	t.Helper()
	path := filepath.Join(dir, "image-mirror.lock.json")
	data, err := json.MarshalIndent(lock{Entries: entries}, "", "  ")
	if err != nil {
		t.Fatalf("marshaling fixture lock file: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("writing fixture lock file: %v", err)
	}
	return path
}

func TestRunResolveOnlyPrintsFreshDigest(t *testing.T) {
	dir := t.TempDir()
	lockFile := writeLockFile(t, dir, lockEntry{
		SourceRepository:  testSourceRepository,
		SourceChannel:     testSourceRepository + ":stable",
		DestinationDigest: "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + strings.Repeat("aa", 32),
	})

	manifest := []byte(`{"mediaType":"single"}`)
	runner := fakeRunner(t, map[string]func() ([]byte, error){
		"skopeo inspect --raw docker://" + testSourceRepository + ":stable": respond(manifest),
	})

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-resolve-only", "-lock-file", lockFile, "-source-repository", testSourceRepository},
		runner, unusedDial(t), unusedImageDigest(t), &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	sum := sha256.Sum256(manifest)
	want := testSourceRepository + "@sha256:" + hex.EncodeToString(sum[:]) + "\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRunResolveOnlyMissingLockEntryFails(t *testing.T) {
	dir := t.TempDir()
	lockFile := writeLockFile(t, dir)

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-resolve-only", "-lock-file", lockFile, "-source-repository", testSourceRepository},
		unusedRunner(t), unusedDial(t), unusedImageDigest(t), &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error for a lock file with no matching entry, got nil")
	}
}

func TestRunReportModeDialError(t *testing.T) {
	dir := t.TempDir()
	lockFile := writeLockFile(t, dir, lockEntry{
		SourceRepository:  testSourceRepository,
		SourceChannel:     testSourceRepository + ":stable",
		DestinationDigest: "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + strings.Repeat("aa", 32),
	})

	wantErr := errors.New("connection refused")
	dial := func(context.Context, dialOptions) (gmp.Connection, error) { return nil, wantErr }

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-lock-file", lockFile, "-source-repository", testSourceRepository},
		unusedRunner(t), dial, unusedImageDigest(t), &stdout, &stderr)
	if !errors.Is(err, wantErr) {
		t.Fatalf("run: got %v, want an error wrapping %v", err, wantErr)
	}
}

func TestRunReportModeDetectsDrift(t *testing.T) {
	dir := t.TempDir()
	pinnedDigest := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + strings.Repeat("aa", 32)
	lockFile := writeLockFile(t, dir, lockEntry{
		SourceRepository:  testSourceRepository,
		SourceChannel:     testSourceRepository + ":stable",
		DestinationDigest: pinnedDigest,
	})

	manifestFile := filepath.Join(dir, "manifest.json")
	writeManifestFixture(t, manifestFile)

	schemaDir := filepath.Join(dir, "schema")
	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		t.Fatalf("mkdir schema dir: %v", err)
	}
	writeGetFixtureSnapshot(t, schemaDir)

	servedDigest := "registry.community.greenbone.net/community/gvmd@sha256:" + strings.Repeat("bb", 32)
	conn := &fakeConn{rawXML: helpResponseXML}
	dial := func(context.Context, dialOptions) (gmp.Connection, error) { return conn, nil }
	imageDigest := func(context.Context, string, string) (string, error) { return servedDigest, nil }

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-lock-file", lockFile, "-source-repository", testSourceRepository,
			"-manifest", manifestFile, "-schema", schemaDir, "-gmp.insecure"},
		unusedRunner(t), dial, imageDigest, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if !conn.closed {
		t.Error("connection was not closed")
	}

	out := stdout.String()
	if !strings.Contains(out, "pinned:  "+pinnedDigest) {
		t.Errorf("stdout = %q, want the pinned digest reported", out)
	}
	if !strings.Contains(out, "served:  "+servedDigest) {
		t.Errorf("stdout = %q, want the served digest reported", out)
	}
	if !strings.Contains(out, "upstream has moved") {
		t.Errorf("stdout = %q, want it to note upstream has moved (pinned and served digests differ)", out)
	}
	if !strings.Contains(out, "Gained") || !strings.Contains(out, "help") {
		t.Errorf("stdout = %q, want get_fixture's sibling command \"help\" reported as gained", out)
	}
}

func TestRunReportModeExpectDigestMismatchFails(t *testing.T) {
	dir := t.TempDir()
	pinnedDigest := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + strings.Repeat("aa", 32)
	lockFile := writeLockFile(t, dir, lockEntry{
		SourceRepository:  testSourceRepository,
		SourceChannel:     testSourceRepository + ":stable",
		DestinationDigest: pinnedDigest,
	})

	manifestFile := filepath.Join(dir, "manifest.json")
	writeManifestFixture(t, manifestFile)
	schemaDir := filepath.Join(dir, "schema")
	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		t.Fatalf("mkdir schema dir: %v", err)
	}
	writeGetFixtureSnapshot(t, schemaDir)

	expectDigest := "registry.community.greenbone.net/community/gvmd@sha256:" + strings.Repeat("cc", 32)
	servedDigest := "registry.community.greenbone.net/community/gvmd@sha256:" + strings.Repeat("bb", 32)
	conn := &fakeConn{rawXML: helpResponseXML}
	dial := func(context.Context, dialOptions) (gmp.Connection, error) { return conn, nil }
	imageDigest := func(context.Context, string, string) (string, error) { return servedDigest, nil }

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-lock-file", lockFile, "-source-repository", testSourceRepository,
			"-manifest", manifestFile, "-schema", schemaDir, "-gmp.insecure",
			"-expect-digest", expectDigest},
		unusedRunner(t), dial, imageDigest, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected an error for a served digest not matching -expect-digest, got nil")
	}
	if !conn.closed {
		t.Error("connection was not closed")
	}
	if stdout.String() != "" {
		t.Errorf("stdout = %q, want no report written on a digest mismatch", stdout.String())
	}
}

func TestRunReportModeNoDrift(t *testing.T) {
	dir := t.TempDir()
	digest := "ghcr.io/brennoo/gvm-mirror/gvmd@sha256:" + strings.Repeat("aa", 32)
	lockFile := writeLockFile(t, dir, lockEntry{
		SourceRepository:  testSourceRepository,
		SourceChannel:     testSourceRepository + ":stable",
		DestinationDigest: digest,
	})

	manifestFile := filepath.Join(dir, "manifest.json")
	writeManifestFixtureBothCommands(t, manifestFile)

	schemaDir := filepath.Join(dir, "schema")
	if err := os.MkdirAll(schemaDir, 0o755); err != nil {
		t.Fatalf("mkdir schema dir: %v", err)
	}
	writeGetFixtureSnapshot(t, schemaDir)
	writeHelpSnapshot(t, schemaDir)

	conn := &fakeConn{rawXML: helpResponseXML}
	dial := func(context.Context, dialOptions) (gmp.Connection, error) { return conn, nil }
	servedDigest := "registry.community.greenbone.net/community/gvmd@sha256:" + strings.Repeat("aa", 32)
	imageDigest := func(context.Context, string, string) (string, error) { return servedDigest, nil }

	var stdout, stderr bytes.Buffer
	err := run(context.Background(),
		[]string{"-lock-file", lockFile, "-source-repository", testSourceRepository,
			"-manifest", manifestFile, "-schema", schemaDir, "-gmp.insecure",
			// A different repository than servedDigest's (upstream vs. this
			// project's GHCR mirror) but the same content hash — proving the
			// -expect-digest check compares hashes, not whole refs.
			"-expect-digest", digest},
		unusedRunner(t), dial, imageDigest, &stdout, &stderr)
	if err != nil {
		t.Fatalf("run: %v", err)
	}

	out := stdout.String()
	if !strings.Contains(out, "upstream has not moved") {
		t.Errorf("stdout = %q, want it to note upstream has not moved (pinned and served digests match)", out)
	}
	if !strings.Contains(out, "No schema drift detected.") {
		t.Errorf("stdout = %q, want no drift reported", out)
	}
}

func unusedRunner(t *testing.T) Runner {
	t.Helper()
	return fakeRunner(t, map[string]func() ([]byte, error){})
}

func unusedDial(t *testing.T) dialFunc {
	t.Helper()
	return func(context.Context, dialOptions) (gmp.Connection, error) {
		t.Fatal("dial should not be called in -resolve-only mode")
		return nil, nil
	}
}

func unusedImageDigest(t *testing.T) imageDigestFunc {
	t.Helper()
	return func(context.Context, string, string) (string, error) {
		t.Fatal("imageDigest should not be called in -resolve-only mode")
		return "", nil
	}
}

// fakeConn implements gmp.Connection with a canned RawXML response, mirroring
// client/client_test.go's mockConn precedent for faking this interface.
// Execute is never called by schemadrift itself — only Authenticate (via
// client.New) inside the real dialLive would call it, and tests inject dial
// directly instead of exercising dialLive.
type fakeConn struct {
	rawXML string
	closed bool
}

func (f *fakeConn) Execute(context.Context, interface{}, interface{}) error { return nil }
func (f *fakeConn) Close() error                                            { f.closed = true; return nil }
func (f *fakeConn) RawXML(context.Context, string) (string, error)          { return f.rawXML, nil }

// helpResponseXML is the same fixture shape as internal/schema/testdata/served/help_response.xml,
// inlined here so this package doesn't reach into another package's testdata.
// It declares two commands, get_fixture and help.
const helpResponseXML = `<help_response status="200" status_text="OK"><schema format="XML" extension="xml" content_type="text/xml"><protocol><version>22.7</version><command><name>get_fixture</name><pattern><e>id</e></pattern><ele><name>id</name><pattern>text</pattern></ele><response><pattern><e>id</e></pattern><ele><name>id</name><pattern>text</pattern></ele></response></command><command><name>help</name><pattern><attrib><name>format</name><pattern>text</pattern></attrib></pattern><response><pattern><e>schema</e></pattern><ele><name>schema</name><pattern>text</pattern></ele></response></command></protocol></schema></help_response>`

func writeManifestFixture(t *testing.T, path string) {
	t.Helper()
	writeManifestJSON(t, path, codegen.Manifest{Families: []codegen.Family{
		{Name: "fixture", Commands: []string{"get_fixture"}},
	}})
}

func writeManifestFixtureBothCommands(t *testing.T, path string) {
	t.Helper()
	writeManifestJSON(t, path, codegen.Manifest{Families: []codegen.Family{
		{Name: "fixture", Commands: []string{"get_fixture"}},
		{Name: "help", Commands: []string{"help"}},
	}})
}

func writeManifestJSON(t *testing.T, path string, m codegen.Manifest) {
	t.Helper()
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshaling fixture manifest: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("writing fixture manifest: %v", err)
	}
}

func elem(name, typ string, kind schema.Kind, card schema.Card, children ...*schema.Node) *schema.Node {
	return &schema.Node{Name: name, Kind: kind, Type: typ, Card: card, Children: children}
}

// writeGetFixtureSnapshot commits get_fixture's shape exactly as helpResponseXML
// parses it (verified against a live schema.LoadServed run), so
// TestRunReportModeNoDrift's committed snapshot and served schema agree.
func writeGetFixtureSnapshot(t *testing.T, dir string) {
	t.Helper()
	cmd := schema.Command{
		Name:     "get_fixture",
		Request:  elem("get_fixture", "", schema.KindElement, schema.CardRequired, elem("id", "text", schema.KindElement, schema.CardRequired)),
		Response: elem("get_fixture_response", "", schema.KindElement, schema.CardRequired, elem("id", "text", schema.KindElement, schema.CardRequired)),
	}
	writeSnapshotCmd(t, dir, cmd)
}

// writeHelpSnapshot commits help's shape exactly as helpResponseXML parses
// it, likewise verified against a live schema.LoadServed run.
func writeHelpSnapshot(t *testing.T, dir string) {
	t.Helper()
	cmd := schema.Command{
		Name:     "help",
		Request:  elem("help", "", schema.KindElement, schema.CardRequired, elem("format", "text", schema.KindAttribute, schema.CardOptional)),
		Response: elem("help_response", "", schema.KindElement, schema.CardRequired, elem("schema", "text", schema.KindElement, schema.CardRequired)),
	}
	writeSnapshotCmd(t, dir, cmd)
}

func writeSnapshotCmd(t *testing.T, dir string, cmd schema.Command) {
	t.Helper()
	data, err := codegen.MarshalSnapshot(cmd)
	if err != nil {
		t.Fatalf("MarshalSnapshot(%s): %v", cmd.Name, err)
	}
	if err := os.WriteFile(codegen.SnapshotPath(dir, cmd.Name), data, 0o600); err != nil {
		t.Fatalf("writing snapshot for %s: %v", cmd.Name, err)
	}
}
