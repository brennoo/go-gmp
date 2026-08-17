package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// commandsDir is the real, committed generated tree, relative to this
// package's directory. Only TestCheckOutputAgreesWithCommittedTree reads it,
// and it never writes to it.
const commandsDir = "../../commands"

// newOutDir materialises a pristine copy of what the generator currently
// produces into a temp directory, plus one hand-written file that carries no
// generated marker. Every mutation test starts from this and breaks exactly
// one thing, so a failure names the category it broke rather than a soup of
// them.
func newOutDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// The adopted families only: this stands in for commands/, and a
	// review-only family's generated output is deliberately not there yet.
	for name, src := range generateAdoptedOrFail(t) {
		writeTestFile(t, filepath.Join(dir, name), src)
	}
	// A hand-written neighbour: commands/ has 118 of these, and none of them
	// may ever be reported obsolete.
	writeTestFile(t, filepath.Join(dir, "handwritten.go"), "package commands\n\ntype HandWritten struct{}\n")
	return dir
}

// writeTestFile writes a fixture into a t.TempDir tree. It deliberately does
// not reuse scratch.go's writeFile: routing manifest-derived command names
// through a non-test helper would pull that helper into gosec's taint graph
// for the sake of a one-line os.WriteFile.
func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("Failed to write %s: %v", path, err)
	}
}

func checkOrFail(t *testing.T, m Manifest, schemaDir, outDir string) Report {
	t.Helper()
	r, err := CheckOutput(m, schemaDir, outDir, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist)
	if err != nil {
		t.Fatalf("CheckOutput failed: %v", err)
	}
	return r
}

// wantOnly asserts that exactly the named category holds exactly want, and
// every other category is empty. Checking the other three is the point: a
// deleted file that also showed up as "obsolete" would be a real bug.
func wantOnly(t *testing.T, r Report, category string, want []string) {
	t.Helper()
	got := map[string][]string{
		"missing":  r.Missing,
		"modified": r.Modified,
		"obsolete": r.Obsolete,
		"snapshot": r.MalformedSnapshots,
	}
	for name, items := range got {
		if name != category {
			if len(items) != 0 {
				t.Errorf("Unexpected %s findings: %v", name, items)
			}
			continue
		}
		if len(items) != len(want) {
			t.Fatalf("%s = %v, want %v", name, items, want)
		}
		for i := range want {
			if items[i] != want[i] {
				t.Errorf("%s[%d] = %q, want %q", name, i, items[i], want[i])
			}
		}
	}
	if r.OK() {
		t.Errorf("Report.OK() = true, want false for a %s finding", category)
	}
}

// TestCheckOutputCleanTree is the baseline the mutation tests are deltas from.
func TestCheckOutputCleanTree(t *testing.T) {
	m := manifestOrFail(t)
	r := checkOrFail(t, m, snapshotDir, newOutDir(t))
	if !r.OK() {
		t.Errorf("A freshly generated tree is not clean:\n%s", r)
	}
	if s := r.String(); s != "" {
		t.Errorf("Report.String() on a clean report = %q, want empty", s)
	}
}

// TestCheckOutputAgreesWithCommittedTree is the regression guard: the whole
// mechanism has to agree with reality as committed right now, using the real
// manifest, the real snapshots and the real commands/ — including its 118
// hand-written files, none of which may be reported obsolete.
func TestCheckOutputAgreesWithCommittedTree(t *testing.T) {
	m := manifestOrFail(t)
	r := checkOrFail(t, m, snapshotDir, commandsDir)
	if !r.OK() {
		t.Errorf("commands/ disagrees with the manifest and schema. Re-run hack/generate.sh.\n%s", r)
	}
}

func TestCheckOutputDetectsMissing(t *testing.T) {
	dir := newOutDir(t)
	gone := filepath.Join(dir, "get_notes.go")
	if err := os.Remove(gone); err != nil {
		t.Fatalf("Failed to remove %s: %v", gone, err)
	}
	wantOnly(t, checkOrFail(t, manifestOrFail(t), snapshotDir, dir), "missing", []string{gone})
}

// TestCheckOutputDetectsMissingSharedFile covers the file the manifest does
// not name a command for: the family's shared dedup output.
func TestCheckOutputDetectsMissingSharedFile(t *testing.T) {
	dir := newOutDir(t)
	fam, ok := manifestOrFail(t).FamilyOf("create_note")
	if !ok {
		t.Fatalf("The manifest no longer declares create_note")
	}
	gone := filepath.Join(dir, fam.SharedFile())
	if err := os.Remove(gone); err != nil {
		t.Fatalf("Failed to remove %s: %v", gone, err)
	}
	wantOnly(t, checkOrFail(t, manifestOrFail(t), snapshotDir, dir), "missing", []string{gone})
}

func TestCheckOutputDetectsModified(t *testing.T) {
	dir := newOutDir(t)
	path := filepath.Join(dir, "create_note.go")
	src, err := os.ReadFile(path) //nolint:gosec // a path this test just created
	if err != nil {
		t.Fatalf("Failed to read %s: %v", path, err)
	}
	// A hand edit that still compiles and still carries the marker: the check
	// has to be a byte comparison, not a marker or parse check.
	edited := strings.Replace(string(src), "type CreateNote struct", "type CreateNote struct // hand-edited", 1)
	if edited == string(src) {
		t.Fatalf("Nothing to edit in %s", path)
	}
	writeTestFile(t, path, edited)
	wantOnly(t, checkOrFail(t, manifestOrFail(t), snapshotDir, dir), "modified", []string{path})
}

func TestCheckOutputDetectsObsolete(t *testing.T) {
	dir := newOutDir(t)
	// A leftover from a family or command that used to be generated: marked,
	// but nothing in the manifest accounts for it. port_range isn't a family
	// in the manifest at all yet, so this can't collide with a real generated
	// file the way "delete_note.go" now would.
	stale := filepath.Join(dir, "delete_port_range.go")
	writeTestFile(t, stale, GeneratedMarker+"\n\npackage commands\n\ntype DeletePortRange struct{}\n")
	wantOnly(t, checkOrFail(t, manifestOrFail(t), snapshotDir, dir), "obsolete", []string{stale})
}

// TestCheckOutputIgnoresUnmarkedAndNonGo proves the obsolete scan keys on the
// marker and the .go suffix, not on "anything the manifest didn't name".
func TestCheckOutputIgnoresUnmarkedAndNonGo(t *testing.T) {
	dir := newOutDir(t)
	writeTestFile(t, filepath.Join(dir, "doc.go"), "// Package commands is hand-written.\npackage commands\n")
	writeTestFile(t, filepath.Join(dir, "notes.json"), GeneratedMarker+"\n")
	sub := filepath.Join(dir, "pagination")
	if err := os.MkdirAll(sub, 0o750); err != nil {
		t.Fatalf("Failed to create %s: %v", sub, err)
	}
	writeTestFile(t, filepath.Join(sub, "page.go"), GeneratedMarker+"\n\npackage pagination\n")

	if r := checkOrFail(t, manifestOrFail(t), snapshotDir, dir); !r.OK() {
		t.Errorf("Unmarked, non-Go and nested files were reported:\n%s", r)
	}
}

func TestCheckOutputDetectsAbsentSnapshot(t *testing.T) {
	// A manifest naming a command with no snapshot file at all, alongside
	// every family newOutDir actually seeded the directory with — otherwise
	// the seeded families' own files would show up as unrelated Obsolete
	// findings, masking the one this test cares about.
	dir := newOutDir(t)
	m := Manifest{Families: append(append([]Family{}, manifestOrFail(t).Adopted().Families...),
		Family{Name: "phantom", Commands: []string{"get_phantoms"}})}

	r := checkOrFail(t, m, snapshotDir, dir)
	if len(r.MalformedSnapshots) != 1 || !strings.HasPrefix(r.MalformedSnapshots[0], "get_phantoms: ") {
		t.Fatalf("MalformedSnapshots = %v, want one get_phantoms entry", r.MalformedSnapshots)
	}
	// The broken family's files are neither expected nor stale: reporting the
	// snapshot once beats reporting it once plus two derived findings.
	if len(r.Missing) != 0 || len(r.Modified) != 0 || len(r.Obsolete) != 0 {
		t.Errorf("A bad snapshot produced derived findings: missing=%v modified=%v obsolete=%v",
			r.Missing, r.Modified, r.Obsolete)
	}
}

func TestCheckOutputDetectsCorruptSnapshot(t *testing.T) {
	dir := newOutDir(t)
	m := manifestOrFail(t)
	schemaCopy := t.TempDir()
	for _, name := range m.AllCommands() {
		src, err := os.ReadFile(SnapshotPath(snapshotDir, name)) //nolint:gosec // a fixed testdata path
		if err != nil {
			t.Fatalf("Failed to read the %s snapshot: %v", name, err)
		}
		writeTestFile(t, SnapshotPath(schemaCopy, name), string(src))
	}
	writeTestFile(t, SnapshotPath(schemaCopy, "modify_note"), "{ this is not JSON")

	r := checkOrFail(t, m, schemaCopy, dir)
	if len(r.MalformedSnapshots) != 1 || !strings.HasPrefix(r.MalformedSnapshots[0], "modify_note: ") {
		t.Fatalf("MalformedSnapshots = %v, want one modify_note entry", r.MalformedSnapshots)
	}
	if r.OK() {
		t.Error("Report.OK() = true despite a corrupt snapshot")
	}
}

// TestCheckOutputWritesNothing is the whole premise of -check: it is a pure
// read of the filesystem. If this ever fails, CI cannot safely run it.
func TestCheckOutputWritesNothing(t *testing.T) {
	dir := newOutDir(t)
	before := dirState(t, dir)

	stale := filepath.Join(dir, "delete_note.go")
	writeTestFile(t, stale, GeneratedMarker+"\n\npackage commands\n")
	if err := os.Remove(filepath.Join(dir, "get_notes.go")); err != nil {
		t.Fatalf("Failed to remove get_notes.go: %v", err)
	}
	dirty := dirState(t, dir)

	if _, err := CheckOutput(manifestOrFail(t), snapshotDir, dir, TypeOverrides, PointerOverrides, RequiredOverrides, FieldOverrides, FieldSuppressions, CommandWarnings, DefectAllowlist); err != nil {
		t.Fatalf("CheckOutput failed: %v", err)
	}
	if got := dirState(t, dir); got != dirty {
		t.Errorf("CheckOutput changed %s:\nbefore %s\nafter  %s", dir, dirty, got)
	}
	if before == dirty {
		t.Fatal("The mutation step did not change the directory; the assertion above proves nothing")
	}
}

// dirState is a name+size summary of dir, enough to catch a created, deleted
// or rewritten file.
func dirState(t *testing.T, dir string) string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read %s: %v", dir, err)
	}
	var b strings.Builder
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			t.Fatalf("Failed to stat %s: %v", e.Name(), err)
		}
		fmt.Fprintf(&b, "%s:%d ", e.Name(), info.Size())
	}
	return b.String()
}
