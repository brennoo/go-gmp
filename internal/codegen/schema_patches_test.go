package codegen

import (
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// withSchemaPatches swaps the package-level SchemaPatches table for the
// duration of one test, restoring the real committed table afterward so
// later tests (and TestSchemaPatchesTableIsValid) never see a leaked
// synthetic entry.
func withSchemaPatches(t *testing.T, patches map[string][]SchemaPatch) {
	t.Helper()
	restore := SchemaPatches
	SchemaPatches = patches
	t.Cleanup(func() { SchemaPatches = restore })
}

// commandWithResultNVTTags returns a synthetic command shaped like
// get_results' own /result/nvt subtree just enough to exercise a patch
// spliced in after "tags" — an optional element, a chardata leaf, nested two
// levels under Response the way result/nvt really sits under
// get_results_response.
func commandWithResultNVTTags() schema.Command {
	c := minimalCommand()
	c.Response.Children = []*schema.Node{
		{
			Name: "result", Kind: schema.KindElement, Card: schema.CardRepeated,
			Children: []*schema.Node{
				{
					Name: "nvt", Kind: schema.KindElement, Card: schema.CardOptional,
					Children: []*schema.Node{
						{Name: "tags", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Chardata: true},
					},
				},
			},
		},
	}
	return c
}

func solutionPatch() []SchemaPatch {
	return []SchemaPatch{{
		Name: "solution", Kind: schema.KindElement, Card: schema.CardOptional,
		Type: "text", Chardata: true, InsertAfter: "tags",
		Reason: "test",
		Children: []SchemaPatch{
			{Name: "type", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text", Reason: "test"},
			{Name: "method", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text", Reason: "test"},
		},
	}}
}

// TestApplySchemaPatchesSplicesIntoGenerate proves the mechanism end to end:
// a patched-in node flows through the real collect/walk/fieldOf/emitStruct
// pipeline into a generated field and struct with the expected wire tags —
// no special-cased emission path of its own.
func TestApplySchemaPatchesSplicesIntoGenerate(t *testing.T) {
	withSchemaPatches(t, map[string][]SchemaPatch{"test_cmd:/result/nvt": solutionPatch()})

	cmd := commandWithResultNVTTags()
	if err := ApplySchemaPatches([]schema.Command{cmd}); err != nil {
		t.Fatalf("ApplySchemaPatches failed: %v", err)
	}

	nvt := cmd.Response.Children[0].Children[0]
	if len(nvt.Children) != 2 || nvt.Children[1].Name != "solution" {
		t.Fatalf("solution was not spliced in after tags: %+v", nvt.Children)
	}

	files, err := Generate([]schema.Command{cmd}, TypeOverrides, PointerOverrides, RequiredOverrides, nil, nil, nil, nil, nil, nil, "", nil, nil, nil)
	if err != nil {
		t.Fatalf("Generate failed on a patched command: %v", err)
	}
	var src string
	for _, s := range files {
		src += s
	}

	for _, want := range []string{
		`xml:"solution,omitempty"`,
		`xml:",chardata"`,
		`xml:"type,attr,omitempty"`,
		`xml:"method,attr,omitempty"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("generated output missing %q; full output:\n%s", want, src)
		}
	}
}

// TestApplySchemaPatchesMissingParentPath proves the fail-closed guarantee
// for a path that resolves in neither the request nor response tree: a typo
// in a patch's key must break the build loudly, not silently no-op.
func TestApplySchemaPatchesMissingParentPath(t *testing.T) {
	withSchemaPatches(t, map[string][]SchemaPatch{"test_cmd:/result/nonexistent": solutionPatch()})

	cmd := commandWithResultNVTTags()
	err := ApplySchemaPatches([]schema.Command{cmd})
	if err == nil || !strings.Contains(err.Error(), "no node at") {
		t.Fatalf("ApplySchemaPatches error = %v, want a missing-path error", err)
	}
}

// TestApplySchemaPatchesWireNameCollision proves the staleness guard: once a
// target already has a child of the same name, applying the patch again (the
// state a live schema catching up would produce) must fail rather than
// silently duplicate or overwrite it.
func TestApplySchemaPatchesWireNameCollision(t *testing.T) {
	withSchemaPatches(t, map[string][]SchemaPatch{"test_cmd:/result/nvt": solutionPatch()})

	cmd := commandWithResultNVTTags()
	nvt := cmd.Response.Children[0].Children[0]
	nvt.Children = append(nvt.Children, &schema.Node{
		Name: "solution", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Chardata: true,
	})

	err := ApplySchemaPatches([]schema.Command{cmd})
	if err == nil || !strings.Contains(err.Error(), "already has a child") {
		t.Fatalf("ApplySchemaPatches error = %v, want a wire-name-collision error", err)
	}
}

// TestValidateSchemaPatchesRejectsBadEntries is a table of every shape
// validateSchemaPatches must reject, mirroring
// TestValidateFieldOverridesRejectsBadEntries for this table.
func TestValidateSchemaPatchesRejectsBadEntries(t *testing.T) {
	tests := []struct {
		name    string
		patches map[string][]SchemaPatch
		wantSub string
	}{
		{
			name:    "key with no colon-slash separator",
			patches: map[string][]SchemaPatch{"test_cmdthing": {{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r"}}},
			wantSub: `key must be "<command>:/<path>"`,
		},
		{
			name:    "empty Name",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{Name: "", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r"}}},
			wantSub: "Name is empty",
		},
		{
			name:    "invalid Kind",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{Name: "n", Kind: schema.Kind(99), Card: schema.CardOptional, Type: "text", Reason: "r"}}},
			wantSub: "neither KindElement nor KindAttribute",
		},
		{
			name:    "invalid Card",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{Name: "n", Kind: schema.KindElement, Card: schema.Card(99), Type: "text", Reason: "r"}}},
			wantSub: "not a valid schema.Card",
		},
		{
			name: "attribute with Children",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{
				Name: "n", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text", Reason: "r",
				Children: []SchemaPatch{{Name: "y", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r"}},
			}}},
			wantSub: "cannot have Children",
		},
		{
			name: "attribute with Chardata",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{
				Name: "n", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text", Chardata: true, Reason: "r",
			}}},
			wantSub: "cannot set Chardata",
		},
		{
			name:    "leaf missing Type",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Reason: "r"}}},
			wantSub: "needs a Type",
		},
		{
			name:    "empty Reason",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: ""}}},
			wantSub: "Reason is empty",
		},
		{
			name: "duplicate name+kind children",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {
				{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r1"},
				{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r2"},
			}},
			wantSub: "duplicate child",
		},
		{
			// findSchemaChild/applySchemaPatchesTo resolve an existing child
			// by name alone, so a same-named element and attribute must be
			// rejected as a pair too — not just an exact {name, kind} repeat.
			name: "duplicate wire name across element and attribute",
			patches: map[string][]SchemaPatch{"test_cmd:/x": {
				{Name: "n", Kind: schema.KindElement, Card: schema.CardOptional, Type: "text", Reason: "r1"},
				{Name: "n", Kind: schema.KindAttribute, Card: schema.CardOptional, Type: "text", Reason: "r2"},
			}},
			wantSub: "duplicate child",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withSchemaPatches(t, tt.patches)
			err := ApplySchemaPatches([]schema.Command{minimalCommand()})
			if err == nil {
				t.Fatalf("ApplySchemaPatches succeeded; want an error containing %q", tt.wantSub)
			}
			if !strings.Contains(err.Error(), tt.wantSub) {
				t.Fatalf("ApplySchemaPatches error = %q, want it to contain %q", err.Error(), tt.wantSub)
			}
		})
	}
}

// TestSchemaPatchesTableIsValid is a regression net: the real, committed
// SchemaPatches table must always pass validateSchemaPatches, so any future
// entry that fails this test has a defect in the entry itself, not in the
// checker.
func TestSchemaPatchesTableIsValid(t *testing.T) {
	if err := validateSchemaPatches(); err != nil {
		t.Fatalf("The committed SchemaPatches table failed validation: %v", err)
	}
}
