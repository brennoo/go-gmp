package codegen

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// writeSnapshot commits cmd to dir the same way TestSnapshotLive does, so
// DetectDrift's LoadSnapshots call has something to read.
func writeSnapshot(t *testing.T, dir string, cmd schema.Command) {
	t.Helper()
	data, err := MarshalSnapshot(cmd)
	if err != nil {
		t.Fatalf("MarshalSnapshot(%s): %v", cmd.Name, err)
	}
	if err := os.WriteFile(SnapshotPath(dir, cmd.Name), data, 0o600); err != nil {
		t.Fatalf("writing snapshot for %s: %v", cmd.Name, err)
	}
}

func elem(name, typ string, card schema.Card, children ...*schema.Node) *schema.Node {
	return &schema.Node{Name: name, Kind: schema.KindElement, Type: typ, Card: card, Children: children}
}

func TestDetectDriftNoChange(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}

	cmd := schema.Command{
		Name:     "get_tasks",
		Request:  elem("get_tasks", "", schema.CardRequired, elem("filter", "string", schema.CardOptional)),
		Response: elem("get_tasks_response", "", schema.CardRequired),
	}
	writeSnapshot(t, dir, cmd)

	report, err := DetectDrift([]schema.Command{cmd}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Gained) != 0 || len(report.Removed) != 0 || len(report.Changed) != 0 {
		t.Fatalf("report = %+v, want no drift", report)
	}
	if got := report.String(); got != "No schema drift detected.\n" {
		t.Errorf("String() = %q, want the no-drift message", got)
	}
}

// TestDetectDriftIgnoresSchemaPatches is the regression the "High" review
// finding on the get_results adoption asked for: DetectDrift's committed-side
// load must stay raw. Before this fix, DetectDrift's LoadSnapshots(snapshotDir,
// ...) call also applied SchemaPatches, so a field that only ever exists in
// the generated Go (never in gvmd's live schema) permanently reported as
// "Removed" — a local generation fix masquerading as upstream schema drift.
func TestDetectDriftIgnoresSchemaPatches(t *testing.T) {
	withSchemaPatches(t, map[string][]SchemaPatch{
		"get_tasks:/get_tasks_response/nvt": {{
			Name: "solution", Kind: schema.KindElement, Card: schema.CardOptional,
			Type: "text", Chardata: true, Reason: "test",
		}},
	})

	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}

	// The committed snapshot on disk is the raw, unpatched shape — exactly
	// like the real testdata/schema/*.json files, which never contain a
	// patched-in node themselves.
	committed := schema.Command{
		Name:     "get_tasks",
		Request:  elem("get_tasks", "", schema.CardRequired),
		Response: elem("get_tasks_response", "", schema.CardRequired, elem("nvt", "", schema.CardOptional)),
	}
	writeSnapshot(t, dir, committed)

	// live is the same raw shape: no gvmd has ever actually declared
	// "solution" here either, so a fresh live extraction is never patched.
	live := schema.Command{
		Name:     "get_tasks",
		Request:  elem("get_tasks", "", schema.CardRequired),
		Response: elem("get_tasks_response", "", schema.CardRequired, elem("nvt", "", schema.CardOptional)),
	}

	report, err := DetectDrift([]schema.Command{live}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Changed) != 0 {
		t.Fatalf("report.Changed = %+v, want no drift: DetectDrift must compare against the raw committed "+
			"snapshot, not one with SchemaPatches applied", report.Changed)
	}
}

func TestDetectDriftGainedAndRemoved(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}

	committed := schema.Command{Name: "get_tasks", Request: elem("get_tasks", "", schema.CardRequired), Response: elem("get_tasks_response", "", schema.CardRequired)}
	writeSnapshot(t, dir, committed)

	live := schema.Command{Name: "get_widgets", Request: elem("get_widgets", "", schema.CardRequired), Response: elem("get_widgets_response", "", schema.CardRequired)}

	report, err := DetectDrift([]schema.Command{live}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Gained) != 1 || report.Gained[0] != "get_widgets" {
		t.Errorf("Gained = %v, want [get_widgets]", report.Gained)
	}
	if len(report.Removed) != 1 || report.Removed[0] != "get_tasks" {
		t.Errorf("Removed = %v, want [get_tasks]", report.Removed)
	}
}

func TestDetectDriftAddedRemovedRetypedNodes(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}

	committed := schema.Command{
		Name: "get_tasks",
		Request: elem("get_tasks", "", schema.CardRequired,
			elem("filter", "string", schema.CardOptional),
			elem("hosts_ordering", "string", schema.CardOptional),
		),
		Response: elem("get_tasks_response", "", schema.CardRequired,
			elem("alive_tests", "string", schema.CardOptional),
		),
	}
	writeSnapshot(t, dir, committed)

	live := schema.Command{
		Name: "get_tasks",
		Request: elem("get_tasks", "", schema.CardRequired,
			elem("filter", "string", schema.CardOptional),
			elem("ignore_pagination", "boolean", schema.CardOptional),
		),
		Response: elem("get_tasks_response", "", schema.CardRequired,
			elem("alive_tests", "string", schema.CardRepeated),
		),
	}

	report, err := DetectDrift([]schema.Command{live}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Changed) != 1 || report.Changed[0].Command != "get_tasks" {
		t.Fatalf("Changed = %+v, want one entry for get_tasks", report.Changed)
	}

	byPath := make(map[string]Drift)
	for _, d := range report.Changed[0].Nodes {
		byPath[d.Path] = d
	}
	if d, ok := byPath["/get_tasks/hosts_ordering"]; !ok || d.Kind != "removed" {
		t.Errorf("hosts_ordering drift = %+v, want removed", d)
	}
	if d, ok := byPath["/get_tasks/ignore_pagination"]; !ok || d.Kind != "added" {
		t.Errorf("ignore_pagination drift = %+v, want added", d)
	}
	if d, ok := byPath["/get_tasks_response/alive_tests"]; !ok || d.Kind != "retyped" || d.Detail != "card: optional -> repeated" {
		t.Errorf("alive_tests drift = %+v, want retyped card optional -> repeated", d)
	}
}

func TestDetectDriftDuplicateSiblingNamesPairByOccurrence(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "filter", Commands: []string{"get_filters"}}}}

	committed := schema.Command{
		Name: "get_filters",
		Request: elem("get_filters", "", schema.CardRequired,
			elem("filters", "string", schema.CardOptional),
			elem("filters", "integer", schema.CardOptional),
		),
		Response: elem("get_filters_response", "", schema.CardRequired),
	}
	writeSnapshot(t, dir, committed)

	// Unchanged shape, same two occurrences in the same order — must not be
	// reported as one removed and one added just because a naive map lookup
	// only kept the last of two same-named children.
	live := schema.Command{
		Name: "get_filters",
		Request: elem("get_filters", "", schema.CardRequired,
			elem("filters", "string", schema.CardOptional),
			elem("filters", "integer", schema.CardOptional),
		),
		Response: elem("get_filters_response", "", schema.CardRequired),
	}

	report, err := DetectDrift([]schema.Command{live}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Changed) != 0 {
		t.Fatalf("Changed = %+v, want no drift for two identically-ordered duplicate siblings", report.Changed)
	}
}

func TestDetectDriftReviewOnlyExcludedFromChanged(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "widget", Commands: []string{"get_widgets"}, ReviewOnly: true}}}

	committed := schema.Command{Name: "get_widgets", Request: elem("get_widgets", "", schema.CardRequired), Response: elem("get_widgets_response", "", schema.CardRequired)}
	writeSnapshot(t, dir, committed)

	live := schema.Command{
		Name:     "get_widgets",
		Request:  elem("get_widgets", "", schema.CardRequired, elem("new_field", "string", schema.CardOptional)),
		Response: elem("get_widgets_response", "", schema.CardRequired),
	}

	report, err := DetectDrift([]schema.Command{live}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Changed) != 0 {
		t.Fatalf("Changed = %+v, want review-only families excluded", report.Changed)
	}
	if len(report.Gained) != 0 || len(report.Removed) != 0 {
		t.Fatalf("report = %+v, want get_widgets counted as present on both sides", report)
	}
}

func TestDetectDriftMissingSnapshotFails(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}
	live := schema.Command{Name: "get_tasks", Request: elem("get_tasks", "", schema.CardRequired), Response: elem("get_tasks_response", "", schema.CardRequired)}

	if _, err := DetectDrift([]schema.Command{live}, m, dir); err == nil {
		t.Fatal("expected an error for a missing committed snapshot, got nil")
	}
}

func TestDetectDriftDuplicateLiveCommandFailsClosed(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{{Name: "task", Commands: []string{"get_tasks"}}}}
	live := schema.Command{Name: "get_tasks", Request: elem("get_tasks", "", schema.CardRequired), Response: elem("get_tasks_response", "", schema.CardRequired)}

	if _, err := DetectDrift([]schema.Command{live, live}, m, dir); err == nil {
		t.Fatal("expected an error for a duplicate live command name, got nil")
	}
}

func TestDetectDriftRanksChangedByNodeCountDescending(t *testing.T) {
	dir := t.TempDir()
	m := Manifest{Families: []Family{
		{Name: "task", Commands: []string{"get_tasks"}},
		{Name: "target", Commands: []string{"get_targets"}},
	}}

	tasksCommitted := schema.Command{Name: "get_tasks", Request: elem("get_tasks", "", schema.CardRequired), Response: elem("get_tasks_response", "", schema.CardRequired)}
	writeSnapshot(t, dir, tasksCommitted)
	targetsCommitted := schema.Command{
		Name: "get_targets",
		Request: elem("get_targets", "", schema.CardRequired,
			elem("a", "string", schema.CardOptional),
			elem("b", "string", schema.CardOptional),
		),
		Response: elem("get_targets_response", "", schema.CardRequired),
	}
	writeSnapshot(t, dir, targetsCommitted)

	tasksLive := schema.Command{
		Name:     "get_tasks",
		Request:  elem("get_tasks", "", schema.CardRequired, elem("new_field", "string", schema.CardOptional)),
		Response: elem("get_tasks_response", "", schema.CardRequired),
	}
	targetsLive := schema.Command{Name: "get_targets", Request: elem("get_targets", "", schema.CardRequired), Response: elem("get_targets_response", "", schema.CardRequired)}

	report, err := DetectDrift([]schema.Command{tasksLive, targetsLive}, m, dir)
	if err != nil {
		t.Fatalf("DetectDrift: %v", err)
	}
	if len(report.Changed) != 2 {
		t.Fatalf("Changed = %+v, want two entries", report.Changed)
	}
	if report.Changed[0].Command != "get_targets" || report.Changed[1].Command != "get_tasks" {
		t.Errorf("Changed order = [%s, %s], want get_targets (2 nodes) before get_tasks (1 node)",
			report.Changed[0].Command, report.Changed[1].Command)
	}

	out := report.String()
	if got := indexOfSub(out, "get_targets"); got == -1 {
		t.Fatalf("String() = %q, want get_targets present", out)
	}
	if got, tasksAt := indexOfSub(out, "get_targets"), indexOfSub(out, "get_tasks"); got >= tasksAt {
		t.Errorf("String() = %q, want get_targets rendered before get_tasks", out)
	}
}

func indexOfSub(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestSnapshotPathRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cmd := schema.Command{Name: "get_tasks", Request: elem("get_tasks", "", schema.CardRequired), Response: elem("get_tasks_response", "", schema.CardRequired)}
	writeSnapshot(t, dir, cmd)
	if _, err := os.Stat(filepath.Join(dir, "get_tasks.json")); err != nil {
		t.Fatalf("snapshot not written where SnapshotPath expects it: %v", err)
	}
}
