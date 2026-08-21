package codegen

import (
	"flag"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// update is this package's regeneration switch, matching the -update idiom
// used elsewhere. It's declared locally, not reused from internal/gmptest,
// to avoid importing that package for a bool and to avoid a flag collision.
var update = flag.Bool("update", false,
	"rewrite the committed schema snapshots and generated goldens")

// TestSnapshot is Stage A's file-based path: an opt-in oracle over
// vendor/XML/commands that never writes and skips when vendor/XML is absent
// (so CI never needs the AGPL input). It proves the fallback path still
// works, but deliberately does not compare against the committed
// snapshots — those are refreshed via TestSnapshotLive and would
// legitimately drift; see TestVendorParity (vendor_parity_test.go) for that
// opt-in comparison.
func TestSnapshot(t *testing.T) {
	cmds := loadSchemaOrSkip(t)

	byName := map[string]schema.Command{}
	for _, c := range cmds {
		byName[c.Name] = c
	}

	for _, name := range manifestOrFail(t).AllCommands() {
		cmd, ok := byName[name]
		if !ok {
			t.Fatalf("%s: not in the parsed schema; %s is stale", name, manifestFile)
		}
		got, err := MarshalSnapshot(cmd)
		if err != nil {
			t.Fatalf("%s: failed to marshal: %v", name, err)
		}
		assertNoProse(t, name, got)
	}
}

// assertNoProse is the licensing guard: it should never fire today since
// schema.Command has no field that can hold a <summary> or <description>,
// but it catches a future field addition upstream before it copies AGPL
// prose into a committed MIT file. It matches JSON keys, not values —
// "description" is a legitimate element name in the schema (e.g.
// get_notes' note/result/description) and may cross freely.
func assertNoProse(t *testing.T, name string, snapshot []byte) {
	t.Helper()
	for _, banned := range []string{"summary", "Summary", "description", "Description", "example", "Example"} {
		if strings.Contains(string(snapshot), `"`+banned+`":`) {
			t.Errorf("%s: snapshot carries a %q key; only names/kinds/cardinalities/types may cross the AGPL boundary",
				name, banned)
		}
	}
}

// TestSnapshotsLoad proves Stage B's input is readable without vendor/XML.
// It is the test that would catch a missing or malformed committed snapshot on
// a clean checkout, where TestSnapshot itself skips.
func TestSnapshotsLoad(t *testing.T) {
	want := manifestOrFail(t).AllCommands()
	cmds, err := LoadSnapshots(snapshotDir, want)
	if err != nil {
		t.Fatalf("Failed to load the committed snapshots: %v", err)
	}
	if len(cmds) != len(want) {
		t.Fatalf("Got %d commands, want %d", len(cmds), len(want))
	}
	for i, c := range cmds {
		if c.Name != want[i] {
			t.Errorf("Snapshot %d is %q, want %q", i, c.Name, want[i])
		}
		if c.Request == nil || c.Response == nil {
			t.Errorf("%s: snapshot is missing a request or response shape", c.Name)
		}
	}
}

func loadSchemaOrSkip(t *testing.T) []schema.Command {
	t.Helper()
	// The vendored tree is a manual copy with no git identity of its own, so
	// there is no revision to supply here today; see schema.Load's doc comment.
	cmds, err := schema.Load(schemaDir, "")
	if err != nil {
		t.Fatalf("Failed to load the schema: %v", err)
	}
	if len(cmds) == 0 {
		t.Skipf("Skipping: no schema at %s (it is gitignored)", schemaDir)
	}
	return cmds
}
