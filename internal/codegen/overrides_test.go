package codegen

import (
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// minimalCommand returns a single-command schema.Command tree just complete
// enough for Generate to run its validation pass: validateFieldOverrides
// checks the FieldOverrides table's own shape unconditionally, without ever
// needing the bad entry's key to match anything this tree declares.
func minimalCommand() schema.Command {
	return schema.Command{
		Name:     "test_cmd",
		Request:  &schema.Node{Name: "test_cmd"},
		Response: &schema.Node{Name: "test_cmd_response"},
	}
}

// TestValidateFieldOverridesRejectsBadEntries is a table of every shape
// validateFieldOverrides must reject, exercised through Generate so the wiring
// from Generate into the check is covered too, not only the check in
// isolation.
func TestValidateFieldOverridesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		fo      map[string][]FieldOverride
		wantSub string
	}{
		{
			name: "empty GoType",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "", Reason: "r"}},
			},
			wantSub: "GoType is empty",
		},
		{
			name: "empty Reason",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "string", Reason: ""}},
			},
			wantSub: "Reason is empty",
		},
		{
			name: "ordinary override with empty Name",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "", GoType: "string", Reason: "r"}},
			},
			wantSub: "Name is empty",
		},
		{
			name: "Chardata combined with Name",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, Name: "n", GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "Chardata combined with Attr",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, Attr: true, GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "Chardata combined with OmitEmpty",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Chardata: true, OmitEmpty: true, GoType: "string", Reason: "r"}},
			},
			wantSub: "mutually exclusive",
		},
		{
			name: "two chardata overrides for one key",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {
					{Chardata: true, GoType: "string", Reason: "r1"},
					{Chardata: true, GoType: "string", Reason: "r2"},
				},
			},
			wantSub: "more than one chardata override",
		},
		{
			name: "duplicate ordinary (Name,Attr) pair for one key",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {
					{Name: "n", GoType: "string", Reason: "r1"},
					{Name: "n", GoType: "int", Reason: "r2"},
				},
			},
			wantSub: "duplicate wire name",
		},
		{
			name: "syntactically invalid GoType func(",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "func(", Reason: "r"}},
			},
			wantSub: "not a valid Go type expression",
		},
		{
			name: "syntactically invalid GoType 1+",
			fo: map[string][]FieldOverride{
				"test_cmd:/x": {{Name: "n", GoType: "1+", Reason: "r"}},
			},
			wantSub: "not a valid Go type expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, tt.fo, nil, nil)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestValidateFieldOverridesAcceptsCommittedTable is a regression net: the
// real, committed FieldOverrides table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateFieldOverridesAcceptsCommittedTable(t *testing.T) {
	g := &gen{fo: FieldOverrides}
	if err := g.validateFieldOverrides(); err != nil {
		t.Fatalf("The committed FieldOverrides table failed validation: %v", err)
	}
}

// commandWithDefects returns minimalCommand() carrying the given defects, so
// checkDefects has something to match against without depending on any real
// committed snapshot.
func commandWithDefects(defects ...schema.Defect) schema.Command {
	c := minimalCommand()
	c.Defects = defects
	return c
}

// TestCheckDefectsRejectsUnallowlisted proves the fail-closed gate: a
// snapshot defect whose {Command, Class, Path} has no DefectAllowlist entry
// must refuse generation, regardless of what Detail says.
func TestCheckDefectsRejectsUnallowlisted(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "pattern declares no content"})
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, nil)
	if err == nil {
		t.Fatal("Generate succeeded with an unallowlisted defect; want an error")
	}
	if !strings.Contains(err.Error(), "not in DefectAllowlist") {
		t.Fatalf("Generate error = %q, want it to mention DefectAllowlist", err)
	}
}

// TestCheckDefectsAcceptsMatchingEntry proves a {Command, Class, Path} match
// permits generation, and that matching ignores Detail: the allowlist key
// carries no Detail field, so two occurrences that differ only in Detail
// still match the same entry.
func TestCheckDefectsAcceptsMatchingEntry(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "some detail text"})
	da := map[DefectKey]DefectAllowlistEntry{
		{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}: {Reason: "test"},
	}
	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, da); err != nil {
		t.Fatalf("Generate failed on an allowlisted defect: %v", err)
	}

	cmd2 := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x", Detail: "totally different detail"})
	if _, err := Generate([]schema.Command{cmd2}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, da); err != nil {
		t.Fatalf("Generate failed on an allowlisted defect with a different Detail: %v", err)
	}
}

// TestCheckDefectsMultipleFailuresAreDeterministic proves that more than one
// unallowlisted defect is reported together, sorted, rather than stopping at
// the first — the same "report every problem in one pass" property
// checkFieldOverridesConsumed already gives FieldOverrides.
func TestCheckDefectsMultipleFailuresAreDeterministic(t *testing.T) {
	cmd := commandWithDefects(
		schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/z"},
		schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/a"},
	)
	_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, nil)
	if err == nil {
		t.Fatal("Generate succeeded with two unallowlisted defects; want an error")
	}
	got := err.Error()
	posA := strings.Index(got, "/test_cmd/a")
	posZ := strings.Index(got, "/test_cmd/z")
	if posA == -1 || posZ == -1 {
		t.Fatalf("Generate error = %q, want both defect paths mentioned", got)
	}
	if posA > posZ {
		t.Fatalf("Generate error = %q, want /test_cmd/a to sort before /test_cmd/z", got)
	}

	second := err
	_, err = Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, nil)
	if err == nil || err.Error() != second.Error() {
		t.Fatalf("Generate error was not deterministic across runs: %q vs %q", second, err)
	}
}

// TestValidateDefectAllowlistRejectsBadEntries is a table of every shape
// validateDefectAllowlist must reject, mirroring
// TestValidateFieldOverridesRejectsBadEntries for the FieldOverrides table.
func TestValidateDefectAllowlistRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		da      map[DefectKey]DefectAllowlistEntry
		wantSub string
	}{
		{
			name:    "empty Command",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "", Class: schema.DefectEmptyPattern, Path: "/x"}: {Reason: "r"}},
			wantSub: "Command is empty",
		},
		{
			name:    "empty Path",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: ""}: {Reason: "r"}},
			wantSub: "Path is empty",
		},
		{
			name:    "unknown Class",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: "not-a-real-class", Path: "/x"}: {Reason: "r"}},
			wantSub: "not a known schema.Defect",
		},
		{
			name:    "empty Reason",
			da:      map[DefectKey]DefectAllowlistEntry{{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/x"}: {Reason: ""}},
			wantSub: "Reason is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, tt.da)
			if err == nil {
				t.Fatalf("Generate succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("Generate error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestCheckDefectAllowlistConsumedRejectsStaleEntry proves the two-sided
// guarantee: an allowlist entry for a command in scope, whose class or path
// no longer matches any generated defect, fails generation just as
// checkFieldOverridesConsumed fails on a stale FieldOverrides entry.
func TestCheckDefectAllowlistConsumedRejectsStaleEntry(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"})

	t.Run("stale path", func(t *testing.T) {
		da := map[DefectKey]DefectAllowlistEntry{
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}:    {Reason: "real"},
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/gone"}: {Reason: "stale"},
		}
		_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, da)
		if err == nil || !strings.Contains(err.Error(), "never matched a generated defect") {
			t.Fatalf("Generate error = %v, want a stale-entry error", err)
		}
	})

	t.Run("stale class", func(t *testing.T) {
		da := map[DefectKey]DefectAllowlistEntry{
			{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}: {Reason: "real"},
			{Command: "test_cmd", Class: schema.DefectNoPattern, Path: "/test_cmd/x"}:    {Reason: "stale"},
		}
		_, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, da)
		if err == nil || !strings.Contains(err.Error(), "never matched a generated defect") {
			t.Fatalf("Generate error = %v, want a stale-entry error", err)
		}
	})
}

// TestCheckDefectAllowlistConsumedIgnoresOutOfScopeEntries proves
// checkDefectAllowlistConsumed's scoping: an allowlist entry for a command
// that is not part of this Generate call is not "stale" for this call, the
// same per-family scoping checkFieldOverridesConsumed already gives
// FieldOverrides — otherwise family-local generation (GenerateAll's whole
// premise) could never pass with a shared allowlist.
func TestCheckDefectAllowlistConsumedIgnoresOutOfScopeEntries(t *testing.T) {
	cmd := commandWithDefects(schema.Defect{Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"})
	da := map[DefectKey]DefectAllowlistEntry{
		{Command: "test_cmd", Class: schema.DefectEmptyPattern, Path: "/test_cmd/x"}:        {Reason: "real"},
		{Command: "some_other_command", Class: schema.DefectEmptyPattern, Path: "/other/y"}: {Reason: "out of scope"},
	}
	if _, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, FieldOverrides, CommandWarnings, da); err != nil {
		t.Fatalf("Generate failed because of an out-of-scope allowlist entry: %v", err)
	}
}

// TestValidateDefectAllowlistAcceptsCommittedTable is a regression net: the
// real, committed DefectAllowlist table must always pass the table-shape
// checks, so any future entry that fails this test has a defect in the entry
// itself, not in the checker.
func TestValidateDefectAllowlistAcceptsCommittedTable(t *testing.T) {
	g := &gen{da: DefectAllowlist}
	if err := g.validateDefectAllowlist(); err != nil {
		t.Fatalf("The committed DefectAllowlist table failed validation: %v", err)
	}
}
