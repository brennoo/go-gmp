package codegen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brennoo/go-gmp/internal/schema"
)

const (
	// schemaDir is gitignored: it is AGPL-licensed and deliberately not
	// vendored into this MIT module. Only Stage A's -update path reads it.
	schemaDir = "../../vendor/XML/commands"
	// snapshotDir holds the committed Stage-A output. Everything downstream
	// of it — Generate, the golden test, CI — reads only this.
	snapshotDir = "testdata/schema"
)

// SnapshotPath is where one command's Stage-A JSON lives under dir.
func SnapshotPath(dir, command string) string {
	return filepath.Join(dir, command+".json")
}

// MarshalSnapshot renders one command as the committed snapshot bytes.
//
// schema.Command is used as the snapshot format directly rather than through a
// parallel IR: its fields are exported, plainly typed and carry exactly the
// names/kinds/cardinalities/type-names that are defensible to move across the
// AGPL boundary. It has no field that could hold <summary> or <example> prose.
func MarshalSnapshot(c schema.Command) ([]byte, error) {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}

// LoadSnapshots reads the committed Stage-A snapshots for the named commands,
// in the order given, exactly as committed — no SchemaPatches applied. This is
// what any raw-schema comparison (DetectDrift, TestVendorParity) must use:
// patching one side of a "does the committed schema match X" comparison but
// not the other would report a local generation fix as upstream schema drift.
//
// The command list is not baked in here; it comes from manifest.json (see
// Manifest.AllCommands and Family.Commands), so there is exactly one checked-in
// declaration of what gets generated.
func LoadSnapshots(dir string, commands []string) ([]schema.Command, error) {
	out := make([]schema.Command, 0, len(commands))
	for _, name := range commands {
		path := SnapshotPath(dir, name)
		data, err := os.ReadFile(path) //nolint:gosec // a fixed testdata path
		if err != nil {
			return nil, fmt.Errorf("%s: %w (run the snapshot test with -update)", name, err)
		}
		var c schema.Command
		if err := json.Unmarshal(data, &c); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		out = append(out, c)
	}
	return out, nil
}

// LoadGenerationModels is LoadSnapshots plus ApplySchemaPatches. Everything
// that turns the schema tree into or validates generated Go — Generate,
// GenerateAll, Validate, and gengmp's write and -check paths — must use this,
// not LoadSnapshots directly, so a SchemaPatches entry actually lands in the
// generated struct.
func LoadGenerationModels(dir string, commands []string) ([]schema.Command, error) {
	out, err := LoadSnapshots(dir, commands)
	if err != nil {
		return nil, err
	}
	if err := ApplySchemaPatches(out); err != nil {
		return nil, err
	}
	return out, nil
}
