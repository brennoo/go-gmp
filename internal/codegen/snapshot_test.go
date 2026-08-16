package codegen

import (
	"flag"
	"strings"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// update is this package's regeneration switch, matching the -update idiom
// internal/schema and internal/gmptest already use. It is declared here rather
// than reused from internal/gmptest because that would make this test binary
// import the fidelity corpus package for a bool, and because two flags of the
// same name only collide inside one binary.
var update = flag.Bool("update", false,
	"rewrite the committed schema snapshots and generated goldens")

// TestSnapshot is Stage A's file-based path: a pure, opt-in comparison
// oracle over vendor/XML/commands, an independently versioned mirror that
// will legitimately drift from whatever gvmd last answered TestSnapshotLive
// (in internal/codegen/snapshot_live_test.go, the primary and only route to
// testdata/schema/*.json since T-16b). It never writes: -update belongs to
// TestSnapshotLive alone, so a file-based run can never overwrite an
// ExtractionServed snapshot with ExtractionFile provenance.
//
// It proves the fallback path itself still works — every manifest command
// parses, marshals, and carries no AGPL prose — but deliberately does NOT
// compare against the committed snapshots: asserting the two independently
// versioned oracles agree is not a regression check, it is a permanent false
// alarm the moment either side moves. See TestVendorParity, in
// vendor_parity_test.go, for that comparison — it is an explicit, opt-in
// drift check (T-17), not part of the default suite.
//
// On a clean checkout vendor/XML is absent and this test skips, which is why
// CI never needs the AGPL input.
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

// assertNoProse is the licensing guard. schema.Command has no field that can
// hold a <summary> or <description>, so this should never fire; it is here so
// that a future field addition upstream fails loudly instead of quietly
// copying AGPL prose into a committed MIT file.
//
// It matches JSON object keys, not values: "description" is a legitimate
// element *name* in the schema (get_notes' note/result/description), and
// carrying that name across is exactly what the boundary permits.
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
