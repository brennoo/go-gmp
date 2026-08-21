package gmptest

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestRequestGolden marshals every curated request in cases_*.go and compares
// the bytes to a committed golden. This is the only test that can catch a
// request-side field name, attribute order, or omitempty regression, since it
// starts from a populated Go value rather than a document already on disk —
// TestRequestRoundTrip can't see a field that should have been emitted but
// silently wasn't.
func TestRequestGolden(t *testing.T) {
	cases := RequestCases()
	if len(cases) == 0 {
		t.Fatal("no request cases registered; see cases_*.go")
	}

	for _, c := range cases {
		t.Run(c.Command+"/"+c.Name, func(t *testing.T) {
			got, err := xml.Marshal(c.Req)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			path := filepath.Join(TestdataDir, requestsDir, c.Command, c.Name+".xml")

			if *Update {
				if err := WriteGolden(path, got); err != nil {
					t.Fatalf("Failed to write the golden: %v", err)
				}
				return
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("No committed golden (%v).\nRun with -update to create it. Got:\n%s", err, got)
			}
			if string(got) != string(want) {
				t.Errorf("Request golden changed.\n--- committed\n%s\n--- now\n%s", want, got)
			}
		})
	}
}

// TestRequestRoundTrip checks, for every fixture that recorded a request
// document, what the library's request type does with it on a round trip.
// Unlike TestRequestGolden, the input is not curated — whatever a gvmd
// changelog example or captured session happened to contain — so this
// measures fidelity against unplanned real-world shapes.
//
// The result rolls up into testdata/REQUEST-FIDELITY.txt rather than each
// fixture's NNN_fidelity.txt, so it can be regenerated and reviewed
// independently of the response-fidelity suite (the same trade
// FIDELITY-SUMMARY.txt already made).
//
// As with TestFidelitySummary, the assertion is "unchanged", not "empty":
// clean round-trips are a minority, and a committed report turns any drift
// into a reviewable diff instead of noise nobody reads.
func TestRequestRoundTrip(t *testing.T) {
	fixtures := loadFixturesOrSkip(t)

	var (
		b                                                strings.Builder
		lines                                            []string
		nTotal, nClean, nDecodeErr, nWithLoss, nAdvisory int
		nSkippedNoRequest                                int
		commandsCovered                                  = map[string]bool{}
		commandsWithLoss                                 = map[string]bool{}
	)

	for _, f := range fixtures {
		if !f.HasRequest() {
			nSkippedNoRequest++
			continue
		}
		cmd, ok := Registry[f.Command]
		if !ok {
			t.Fatalf("No registry entry for %q", f.Command)
		}
		doc, err := f.Request()
		if err != nil {
			t.Fatalf("Failed to read %s: %v", f.RequestPath(), err)
		}

		fid := Compare(doc, cmd.NewReq())
		nTotal++
		commandsCovered[f.Command] = true

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
			nWithLoss++
			commandsWithLoss[f.Command] = true
		}

		lines = append(lines, fmt.Sprintf("%-28s %-9s %-3s %-13s dropped=%-3d truncated=%-3d changed=%-3d added=%d",
			f.Command, f.Provenance, f.ID, status, len(fid.Dropped), len(fid.Truncated), len(fid.Changed), len(fid.Added)))
	}

	sort.Strings(lines)

	b.WriteString("# Request fidelity summary\n#\n")
	b.WriteString("# What each recorded *_request.xml loses when it round-trips through the\n")
	b.WriteString("# library's request type via Registry[command].NewReq. Complements\n")
	b.WriteString("# FIDELITY-SUMMARY.txt (which does the same for responses) and\n")
	b.WriteString("# TestRequestGolden (which starts from curated Go values instead of\n")
	b.WriteString("# recorded documents, and is the one that can catch a field the library\n")
	b.WriteString("# silently fails to ever emit).\n")
	b.WriteString("# Regenerate with: go test ./internal/gmptest -run TestRequestRoundTrip -update\n\n")
	fmt.Fprintf(&b, "requests-compared: %d   clean: %d   with-loss: %d   decode-errors: %d   advisory: %d   skipped-no-request: %d\n",
		nTotal, nClean, nWithLoss, nDecodeErr, nAdvisory, nSkippedNoRequest)
	fmt.Fprintf(&b, "commands with a recorded request: %d of %d   commands with loss: %d\n\n",
		len(commandsCovered), len(Registry), len(commandsWithLoss))

	b.WriteString("## Per fixture\n\n")
	for _, l := range lines {
		b.WriteString(l + "\n")
	}

	got := b.String()
	path := filepath.Join(TestdataDir, "REQUEST-FIDELITY.txt")

	if *Update {
		if err := WriteGolden(path, []byte(got)); err != nil {
			t.Fatalf("Failed to write the summary: %v", err)
		}
		t.Logf("requests-compared=%d clean=%d with-loss=%d decode-errors=%d advisory=%d commands-covered=%d",
			nTotal, nClean, nWithLoss, nDecodeErr, nAdvisory, len(commandsCovered))
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("No committed summary (%v). Run with -update.", err)
	}
	if got != string(want) {
		t.Errorf("Request fidelity summary changed. Run with -update to refresh.\n--- committed\n%s\n--- now\n%s", want, got)
	}
}
