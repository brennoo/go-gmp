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

// families is the prototype's scope: one small, representative command family
// (T-9 recommends proving the pipeline on one family before the other 119).
// Hardcoded on purpose; there is no flag to widen it until the output of this
// slice has been reviewed and swapped in.
var families = []string{"create_note", "get_notes", "modify_note"}

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

// LoadSnapshots reads the committed Stage-A snapshots for every command in
// families, in that (sorted) order. This is Stage B's only input: it never
// touches vendor/XML.
func LoadSnapshots(dir string) ([]schema.Command, error) {
	out := make([]schema.Command, 0, len(families))
	for _, name := range families {
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
