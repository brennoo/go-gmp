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
			_, err := Generate([]schema.Command{minimalCommand()}, TypeOverrides, PointerOverrides, tt.fo, nil)
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
