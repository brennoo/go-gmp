package codegen

import (
	"encoding/json"
	"flag"
	"os"
	"testing"

	"github.com/brennoo/go-gmp/internal/schema"
)

// checkVendorParity opts into TestVendorParity; see its doc comment.
var checkVendorParity = flag.Bool("vendor-parity", false,
	"compare a local vendor/XML/commands mirror against the committed, live-derived "+
		"testdata/schema/*.json — a manual drift check (T-17), not a CI gate")

// TestVendorParity is the manual counterpart to T-16b's now-closed parity
// check between file-based extraction (vendor/XML/commands, via TestSnapshot)
// and served extraction (a live gvmd, via TestSnapshotLive) — the two sides
// of Stage A. It is skipped by default and only runs under -vendor-parity: a
// local vendor mirror and the gvmd that last refreshed testdata/schema/*.json
// are two independently versioned oracles that will legitimately disagree as
// either one moves, so failures here are informational drift reports for a
// human (exactly what T-17 wants automated one day), not a correctness bug
// and not something the default `go test ./...` gate may depend on.
//
// A failure here means: either the local vendor/XML/commands mirror is
// stale relative to the gvmd that produced the committed snapshots, or that
// gvmd has drifted further since — not that generated code is wrong.
func TestVendorParity(t *testing.T) {
	if !*checkVendorParity {
		t.Skip("run with -vendor-parity to diff a local vendor/XML/commands mirror against the committed served snapshots")
	}

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

		path := SnapshotPath(snapshotDir, name)
		want, readErr := os.ReadFile(path) //nolint:gosec // a fixed testdata path
		if readErr != nil {
			t.Fatalf("No committed snapshot for %s (%v).", name, readErr)
		}
		var wantCmd schema.Command
		if err := json.Unmarshal(want, &wantCmd); err != nil {
			t.Fatalf("%s: failed to parse the committed snapshot: %v", name, err)
		}

		// File-based extraction always produces file Provenance; the
		// committed snapshots carry served Provenance and no File.
		// Comparing shape only — File and Provenance excluded — is what
		// makes this a schema-drift check instead of the two routes'
		// provenance disagreeing on every run.
		gotShape, err := MarshalSnapshot(schema.Command{Name: cmd.Name, Request: cmd.Request, Response: cmd.Response, Defects: cmd.Defects})
		if err != nil {
			t.Fatalf("%s: failed to marshal: %v", name, err)
		}
		wantShape, err := MarshalSnapshot(schema.Command{Name: wantCmd.Name, Request: wantCmd.Request, Response: wantCmd.Response, Defects: wantCmd.Defects})
		if err != nil {
			t.Fatalf("%s: failed to marshal: %v", name, err)
		}
		if string(gotShape) != string(wantShape) {
			t.Errorf("%s: local vendor mirror disagrees with the committed served snapshot.\n"+
				"committed %d bytes, local vendor mirror %d bytes", name, len(wantShape), len(gotShape))
		}
	}
}
