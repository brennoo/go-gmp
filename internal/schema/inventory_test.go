// Package schema_test holds the optional-scalar inventory test. It must be an
// external test package: internal/codegen imports internal/schema, so an
// internal (package schema) test file importing codegen would be an import
// cycle. codegen's committed snapshots are the always-available full-corpus
// source this golden needs — see inventory_test.go's doc comment.
package schema_test

import (
	"os"
	"testing"

	"github.com/brennoo/go-gmp/internal/codegen"
	"github.com/brennoo/go-gmp/internal/gmptest"
	"github.com/brennoo/go-gmp/internal/schema"
)

// inventoryPath is the committed golden. Unlike CONFORMANCE.md, this must run
// unconditionally in CI: it is sourced from internal/codegen's committed,
// MIT-licensed snapshots rather than the gitignored AGPL vendor tree, so
// there is nothing here for it to skip on.
const inventoryPath = "testdata/OPTIONAL_SCALAR_INVENTORY.md"

// TestOptionalScalarInventory builds the inventory and asserts it against the
// committed golden. Any generator change that flips a scalar field's
// pointer-ness, omitempty, or adds/removes an optional scalar changes this
// file and fails here until reviewed and -update'd.
func TestOptionalScalarInventory(t *testing.T) {
	inv := buildInventory(t)
	got := inv.String()

	if *gmptest.Update {
		if err := os.WriteFile(inventoryPath, []byte(got), 0o600); err != nil {
			t.Fatalf("Failed to write the inventory: %v", err)
		}
		return
	}

	want, err := os.ReadFile(inventoryPath)
	if err != nil {
		t.Fatalf("No committed inventory (%v).\nRun with -update to create it.", err)
	}
	if got != string(want) {
		t.Errorf("Optional scalar inventory changed. Re-run with -update and review the diff.\n"+
			"committed %d bytes, now %d bytes", len(want), len(got))
	}
}

// TestOptionalScalarInventoryIsDeterministic guards the golden's only real
// precondition, mirroring TestReportIsDeterministic.
func TestOptionalScalarInventoryIsDeterministic(t *testing.T) {
	a := buildInventory(t).String()
	b := buildInventory(t).String()
	if a != b {
		t.Errorf("Two runs produced different inventories; the golden cannot be stable")
	}
}

func buildInventory(t *testing.T) schema.Inventory {
	t.Helper()

	m, err := codegen.LoadManifest("../codegen/manifest.json")
	if err != nil {
		t.Fatalf("Failed to load the manifest: %v", err)
	}
	commandNames := m.Adopted().AllCommands()

	cmds, err := codegen.LoadGenerationModels("../codegen/testdata/schema", commandNames)
	if err != nil {
		t.Fatalf("Failed to load generation models: %v", err)
	}

	var inv schema.Inventory
	c := &schema.Comparer{OnOptionalScalar: func(cmd, side string, s *schema.Node, g *schema.GoNode, path string) {
		inv.Rows = append(inv.Rows, schema.ScalarRow{
			Command:    cmd,
			Side:       side,
			Path:       path,
			SchemaType: s.Type,
			Alts:       s.Alts,
			GoType:     g.GoType,
			Pointer:    g.Pointer,
			OmitEmpty:  g.OmitEmpty,
		})
	}}
	for _, cmd := range cmds {
		entry, ok := gmptest.Registry[cmd.Name]
		if !ok {
			t.Fatalf("%s: in the manifest but not in gmptest.Registry", cmd.Name)
		}
		c.Compare(cmd, schema.ModelOf(entry.NewReq()), schema.ModelOf(entry.NewResp()))
	}
	return inv
}
