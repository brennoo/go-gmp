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
// the bytes to a committed golden.
//
// This is the only test that can catch a request-side field name, attribute
// order, or omitempty regression: it starts from a populated Go value and
// checks what encoding/xml actually produces. TestRequestRoundTrip cannot see
// this — it starts from a document already on disk, so a field that should
// have been emitted but silently wasn't never shows up as a source path to
// begin with.
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
//
// Unlike TestRequestGolden, the input here is not curated: it is whatever a
// gvmd changelog example or a captured session happened to contain, so this
// test measures fidelity against unplanned real-world shapes rather than the
// cases this package's authors thought to write.
//
// The result is rolled up into one file, testdata/REQUEST-FIDELITY.txt, rather
// than appended into each fixture's NNN_fidelity.txt. Two independent tests
// writing into the same file under -update would race on which one's content
// survives depending on run order; a dedicated file makes the two suites
// (response fidelity, request fidelity) independently regenerable and
// independently reviewable, at the cost of one more allowed name in
// allowedUpdatePath — the same trade the existing FIDELITY-SUMMARY.txt already
// made for the response side.
//
// As with TestFidelitySummary, the assertion is "unchanged", not "empty":
// most of these fixtures exist to document upstream examples, and clean
// round-trips are a minority. A committed report turns any drift, in either
// direction, into a reviewable diff instead of noise nobody reads.
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

		lines = append(lines, fmt.Sprintf("%-28s %-9s %-3s %-13s dropped=%-3d truncated=%-3d changed=%-3d added=%-3d",
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
