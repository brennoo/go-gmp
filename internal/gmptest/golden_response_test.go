package gmptest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestResponseFidelity checks every fixture's round-trip report against a
// committed golden.
//
// The assertion is that the report is *unchanged*, not that it is empty. An
// empty report is unreachable: 103 bool fields tagged omitempty cannot
// round-trip a false value, so they will be listed until T-5 converts them to
// *bool. Committing the report instead makes the volume irrelevant and turns
// every genuine fix into a reviewable shrinkage diff, while a server that
// starts sending a new element shows up as a single added line.
func TestResponseFidelity(t *testing.T) {
	fixtures := loadFixturesOrSkip(t)

	for _, f := range fixtures {
		t.Run(string(f.Provenance)+"/"+f.Command+"/"+f.ID, func(t *testing.T) {
			cmd, ok := Registry[f.Command]
			if !ok {
				t.Fatalf("No registry entry for %q", f.Command)
			}
			doc, err := f.Response()
			if err != nil {
				t.Fatalf("Failed to read the response: %v", err)
			}

			got := Compare(doc, cmd.NewResp()).Report(fidelityHeader(f))

			if *Update {
				if err := WriteGolden(f.FidelityPath(), []byte(got)); err != nil {
					t.Fatalf("Failed to write the golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(f.FidelityPath())
			if err != nil {
				t.Fatalf("No committed fidelity report (%v).\nRun with -update to create it. Got:\n%s", err, got)
			}
			if got != string(want) {
				t.Errorf("Fidelity report changed.\n--- committed\n%s\n--- now\n%s", want, got)
			}
		})
	}
}

func fidelityHeader(f Fixture) string {
	return fmt.Sprintf("# fidelity: %s / %s / %s\n# provenance: %s\n", f.Command, f.Provenance, f.ID, f.Meta.Source)
}

// TestFidelitySummary maintains the roll-up. It is the file worth actually
// reading: one line per fixture, so a regression anywhere is visible without
// opening 154 reports.
func TestFidelitySummary(t *testing.T) {
	fixtures := loadFixturesOrSkip(t)

	var (
		b                                       strings.Builder
		lines                                   []string
		nClean, nDecodeErr, nDropped, nTotal    int
		droppedPaths                            = map[string]int{}
		commandsWithLoss                        = map[string]bool{}
		totalDropped, totalTruncated, totalAdds int
		nAdvisory                               int
	)

	for _, f := range fixtures {
		cmd, ok := Registry[f.Command]
		if !ok {
			continue
		}
		doc, err := f.Response()
		if err != nil {
			t.Fatalf("Failed to read %s: %v", f.ResponsePath(), err)
		}
		fid := Compare(doc, cmd.NewResp())
		nTotal++

		status := "ok"
		switch {
		case fid.DecodeErr != nil:
			status = "DECODE-ERROR"
			nDecodeErr++
		case fid.Advisory():
			status = "ADVISORY"
			nAdvisory++
		case fid.Clean():
			status = "clean"
			nClean++
		default:
			nDropped++
			commandsWithLoss[f.Command] = true
		}
		totalDropped += len(fid.Dropped)
		totalTruncated += len(fid.Truncated)
		totalAdds += len(fid.Added)
		for _, d := range fid.Dropped {
			droppedPaths[d.Path]++
		}

		lines = append(lines, fmt.Sprintf("%-28s %-9s %-3s %-13s dropped=%-3d truncated=%-3d changed=%-3d added=%d",
			f.Command, f.Provenance, f.ID, status, len(fid.Dropped), len(fid.Truncated), len(fid.Changed), len(fid.Added)))
	}

	sort.Strings(lines)

	b.WriteString("# Protocol fidelity summary\n#\n")
	b.WriteString("# What each response document loses when it round-trips through its Go type.\n")
	b.WriteString("# Regenerate with: go test ./internal/gmptest -run TestFidelitySummary -update\n")
	b.WriteString("#\n# 'dropped' is the count of distinct paths no Go field captures; a dropped\n")
	b.WriteString("# parent folds its descendants, so one line can stand for a whole subtree.\n")
	b.WriteString("# 'changed' is a path that survived carrying a different value, e.g. a GMP\n")
	b.WriteString("# boolean '1' re-emitted as 'true'. Structural retention is not fidelity.\n")
	b.WriteString("# 'added' is elements the round trip invents (non-pointer fields without\n")
	b.WriteString("# omitempty, zero time.Time) and is not loss.\n\n")
	fmt.Fprintf(&b, "fixtures: %d   clean: %d   with-loss: %d   decode-errors: %d   advisory: %d\n",
		nTotal, nClean, nDropped, nDecodeErr, nAdvisory)
	fmt.Fprintf(&b, "commands with loss: %d of %d\n", len(commandsWithLoss), len(Registry))
	fmt.Fprintf(&b, "distinct dropped paths: %d\n\n", len(droppedPaths))

	b.WriteString("## Per fixture\n\n")
	for _, l := range lines {
		b.WriteString(l + "\n")
	}

	got := b.String()
	path := filepath.Join(TestdataDir, "FIDELITY-SUMMARY.txt")

	if *Update {
		if err := WriteGolden(path, []byte(got)); err != nil {
			t.Fatalf("Failed to write the summary: %v", err)
		}
		t.Logf("fixtures=%d clean=%d with-loss=%d decode-errors=%d advisory=%d distinct-dropped-paths=%d",
			nTotal, nClean, nDropped, nDecodeErr, nAdvisory, len(droppedPaths))
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("No committed summary (%v). Run with -update.", err)
	}
	if got != string(want) {
		t.Errorf("Fidelity summary changed. Run with -update to refresh.\n--- committed\n%s\n--- now\n%s", want, got)
	}
}

func loadFixturesOrSkip(t *testing.T) []Fixture {
	t.Helper()
	fixtures, err := LoadFixtures(TestdataDir)
	if err != nil {
		t.Fatalf("Failed to load fixtures: %v", err)
	}
	if len(fixtures) == 0 {
		t.Skip("Skipping: no fixtures present")
	}
	return fixtures
}
