package schema

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/brennoo/go-gmp/internal/gmptest"
)

// update is internal/gmptest's -update flag, reused rather than redeclared:
// registering a second flag of the same name in one test binary panics at
// init, and one switch for every derived golden in the repository is the
// behaviour to want anyway.
var update = gmptest.Update

const (
	// schemaDir is gitignored: it is AGPL-licensed and deliberately not
	// vendored into this MIT module. When it is absent the report cannot be
	// regenerated, and the test skips rather than fails.
	schemaDir = "../../vendor/XML/commands"
	// fidelityDir is read, never written. Another tool owns it.
	fidelityDir = "../gmptest/testdata"
	reportPath  = "testdata/CONFORMANCE.md"
)

// TestConformanceReport builds the report and asserts it against the committed
// golden.
//
// The assertion is "unchanged", not "empty". 122 commands were transcribed by
// hand and the report says what that cost; turning every finding into a test
// failure would mean either a permanently red suite or a suppression list
// nobody reads. A committed report makes any drift, in either direction, a
// reviewable diff.
func TestConformanceReport(t *testing.T) {
	rep := buildReport(t)
	got := rep.String()

	if *update {
		if err := os.MkdirAll(filepath.Dir(reportPath), 0o750); err != nil {
			t.Fatalf("Failed to create testdata: %v", err)
		}
		if err := os.WriteFile(reportPath, []byte(got), 0o600); err != nil {
			t.Fatalf("Failed to write the report: %v", err)
		}
		return
	}

	want, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("No committed report (%v).\nRun with -update to create it.", err)
	}
	if got != string(want) {
		t.Errorf("Conformance report changed. Re-run with -update and review the diff.\n"+
			"committed %d bytes, now %d bytes", len(want), len(got))
	}
}

// TestReportIsDeterministic guards the golden's only real precondition.
func TestReportIsDeterministic(t *testing.T) {
	a := buildReport(t).String()
	b := buildReport(t).String()
	if a != b {
		t.Errorf("Two runs produced different reports; the golden cannot be stable")
	}
}

// TestEveryCommandParses is the parser's smoke test against the real corpus.
// A file that fails to parse would silently drop a command from the report.
func TestEveryCommandParses(t *testing.T) {
	cmds := loadSchemaOrSkip(t)
	if len(cmds) < 100 {
		t.Fatalf("Got %d commands, expected the full corpus", len(cmds))
	}
	for _, c := range cmds {
		if c.Name == "" {
			t.Errorf("%s: parsed a command with no name", c.File)
		}
		if c.Request == nil || c.Response == nil {
			t.Errorf("%s: missing a request or response shape", c.Name)
		}
	}
}

// TestGoModelMatchesRegistry checks the reflection side independently of the
// schema, so a bug there cannot be mistaken for a transcription error.
func TestGoModelMatchesRegistry(t *testing.T) {
	for name, cmd := range gmptest.Registry {
		req := ModelOf(cmd.NewReq())
		if req.Root == nil || req.Root.Name != name {
			t.Errorf("%s: request root is %q, want %q", name, rootNameOf(req), name)
		}
		resp := ModelOf(cmd.NewResp())
		if resp.Root == nil || resp.Root.Name != name+"_response" {
			t.Errorf("%s: response root is %q, want %q_response", name, rootNameOf(resp), name)
		}
	}
}

func rootNameOf(r GoRoot) string {
	if r.Root == nil {
		return "<nil>"
	}
	return r.Root.Name
}

// TestLoadWithoutSchemaIsNotAnError is what makes the skip in
// loadSchemaOrSkip correct rather than lucky: vendor/ is gitignored, so on a
// clean checkout this is the path every run takes.
func TestLoadWithoutSchemaIsNotAnError(t *testing.T) {
	cmds, err := Load(filepath.Join(t.TempDir(), "absent"), "")
	if err != nil {
		t.Fatalf("A missing schema directory must not be an error, got %v", err)
	}
	if len(cmds) != 0 {
		t.Fatalf("Got %d commands from an empty directory", len(cmds))
	}
	dropped, err := LoadDroppedPaths(filepath.Join(t.TempDir(), "absent"))
	if err != nil || len(dropped) != 0 {
		t.Fatalf("Got %v, %v; want an empty set and no error", dropped, err)
	}
}

func buildReport(t *testing.T) Report {
	t.Helper()
	cmds := loadSchemaOrSkip(t)

	dropped, err := LoadDroppedPaths(fidelityDir)
	if err != nil {
		t.Fatalf("Failed to read the fidelity reports: %v", err)
	}
	c := &Comparer{Dropped: dropped}

	rep := Report{SchemaDir: "vendor/XML/commands"}
	inSchema := map[string]bool{}
	for _, cmd := range cmds {
		inSchema[cmd.Name] = true
		entry, ok := gmptest.Registry[cmd.Name]
		if !ok {
			rep.SchemaOnly = append(rep.SchemaOnly, cmd.Name)
			continue
		}
		rep.Results = append(rep.Results, c.Compare(cmd, ModelOf(entry.NewReq()), ModelOf(entry.NewResp())))
	}
	for name := range gmptest.Registry {
		if !inSchema[name] {
			rep.GoOnly = append(rep.GoOnly, name)
		}
	}
	sort.Strings(rep.SchemaOnly)
	sort.Strings(rep.GoOnly)
	sort.Slice(rep.Results, func(i, j int) bool { return rep.Results[i].Command < rep.Results[j].Command })
	return rep
}

func loadSchemaOrSkip(t *testing.T) []Command {
	t.Helper()
	// The vendored tree is a manual copy with no git identity of its own, so
	// there is no revision to supply here today; see Load's doc comment.
	cmds, err := Load(schemaDir, "")
	if err != nil {
		t.Fatalf("Failed to load the schema: %v", err)
	}
	if len(cmds) == 0 {
		t.Skipf("Skipping: no schema at %s (it is gitignored)", schemaDir)
	}
	return cmds
}
