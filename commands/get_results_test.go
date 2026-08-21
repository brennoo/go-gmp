package commands

import (
	"encoding/xml"
	"strconv"
	"strings"
	"testing"
)

// The served schema declares result/nvt/severities/severity as repeated;
// pin that both entries survive the round trip rather than only the last one
// winning, as a scalar field would.
func TestGetResultsNVTSeveritiesRepeats(t *testing.T) {
	in := `<nvt><severities score="7.5">` +
		`<severity type="cvss_base_v2"><origin>NVD</origin><date>2020-01-01T00:00:00Z</date><score>75</score><value>AV:N/AC:L</value></severity>` +
		`<severity type="cvss_base_v3"><origin>NVD</origin><date>2020-01-02T00:00:00Z</date><score>78</score><value>AV:N/AC:L/PR:N/UI:N</value></severity>` +
		`</severities></nvt>`
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(in), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(nvt.Severities.Severity) != 2 {
		t.Fatalf("Severity = %+v, want 2 entries", nvt.Severities.Severity)
	}
	if nvt.Severities.Severity[0].Type != "cvss_base_v2" || nvt.Severities.Severity[1].Type != "cvss_base_v3" {
		t.Fatalf("Severity = %+v, want both types preserved in order", nvt.Severities.Severity)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := strings.Count(string(out), "<severity "); got != 2 {
		t.Fatalf("Marshal = %s, want 2 <severity> elements, got %d", out, got)
	}
}

// Upstream fixture 001 sends cvss_base as a real decimal ("5.0") even though
// the served schema declares it integer; a plain int would truncate this on
// round trip, which is why it's typed float64.
func TestGetResultsNVTCVSSBaseDecimal(t *testing.T) {
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(`<nvt><cvss_base>5.0</cvss_base></nvt>`), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.CVSSBase != 5.0 {
		t.Fatalf("CVSSBase = %v, want 5.0", nvt.CVSSBase)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	start := strings.Index(string(out), "<cvss_base>") + len("<cvss_base>")
	end := strings.Index(string(out), "</cvss_base>")
	got, err := strconv.ParseFloat(string(out)[start:end], 64)
	if err != nil || got != 5.0 {
		t.Fatalf("Marshal = %s, want cvss_base to round-trip as 5.0", out)
	}
}

// All 4 live get_results fixtures send <solution type=""></solution>: present
// but empty. Solution must decode to a non-nil pointer, not be dropped as
// though the element were absent.
func TestGetResultsNVTSolutionEmptyPresent(t *testing.T) {
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(`<nvt><solution type=''></solution></nvt>`), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.Solution == nil {
		t.Fatal("Solution = nil, want a non-nil pointer for a present-but-empty element")
	}
	if nvt.Solution.Value != "" || nvt.Solution.Type != "" || nvt.Solution.Method != "" {
		t.Fatalf("Solution = %+v, want all-empty fields", nvt.Solution)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), "<solution") {
		t.Fatalf("Marshal = %s, want the solution element preserved despite all-zero fields", out)
	}
	// Accepted, documented loss (docs/review-2026-08.md's 2026-08-19
	// addendum): Type's `omitempty` can't distinguish explicit-empty from
	// absent, so type="" never survives remarshal. Pin the baseline instead
	// of letting this test read as if the whole observed shape round-trips.
	if strings.Contains(string(out), "type=") {
		t.Fatalf("Marshal = %s, want type=\"\" dropped (accepted omitempty loss), not preserved", out)
	}
}

// TestGetResultsNVTSolutionNonEmptyTypeSurvives is the control for the loss
// TestGetResultsNVTSolutionEmptyPresent pins: omitempty only drops an
// explicit-empty attribute, not a genuinely present one, so a non-empty
// type must still round-trip.
func TestGetResultsNVTSolutionNonEmptyTypeSurvives(t *testing.T) {
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(`<nvt><solution type="Workaround">apply the patch</solution></nvt>`), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.Solution == nil || nvt.Solution.Type != "Workaround" || nvt.Solution.Value != "apply the patch" {
		t.Fatalf("Solution = %+v, want Type=Workaround, Value=%q", nvt.Solution, "apply the patch")
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), `type="Workaround"`) {
		t.Fatalf("Marshal = %s, want a non-empty type attribute preserved", out)
	}
}

func TestGetResultsOriginalSeverityPresent(t *testing.T) {
	var r Result
	if err := xml.Unmarshal([]byte(`<result><original_severity>7.5</original_severity></result>`), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r.OriginalSeverity == nil || *r.OriginalSeverity != 7.5 {
		t.Fatalf("OriginalSeverity = %v, want a pointer to 7.5", r.OriginalSeverity)
	}

	out, err := xml.Marshal(&r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), "<original_severity>7.5</original_severity>") {
		t.Fatalf("Marshal = %s, want original_severity preserved", out)
	}
}

func TestGetResultsOriginalSeverityAbsent(t *testing.T) {
	var r Result
	if err := xml.Unmarshal([]byte(`<result></result>`), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r.OriginalSeverity != nil {
		t.Fatalf("OriginalSeverity = %v, want nil since the source never sent it", r.OriginalSeverity)
	}

	out, err := xml.Marshal(&r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(out), "original_severity") {
		t.Fatalf("Marshal = %s, want no fabricated original_severity", out)
	}
}

// T-26: nvt/refs is schema-required, so a plain (non-pointer) field would
// fabricate an empty <refs></refs> on remarshal for every NVT that never had
// one — exactly the fabrication live/get_results/001_response.xml showed (1
// of 5 <nvt> elements has <refs>, the other 4 don't; decoding then remarshaling
// with the pre-fix non-pointer field produced 4 fabricated empty <refs>).
func TestGetResultsNVTRefsAbsentStaysAbsent(t *testing.T) {
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(`<nvt></nvt>`), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.Refs != nil {
		t.Fatalf("Refs = %+v, want nil since the source never sent <refs>", nvt.Refs)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if strings.Contains(string(out), "<refs") {
		t.Fatalf("Marshal = %s, want no fabricated <refs> element", out)
	}
}

func TestGetResultsNVTRefsPresentRoundTrips(t *testing.T) {
	var nvt ResultNVT
	in := `<nvt><refs><ref id="CVE-2020-0001" type="cve"></ref></refs></nvt>`
	if err := xml.Unmarshal([]byte(in), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.Refs == nil || len(nvt.Refs.Ref) != 1 || nvt.Refs.Ref[0].ID != "CVE-2020-0001" {
		t.Fatalf("Refs = %+v, want one ref with ID=CVE-2020-0001", nvt.Refs)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), `<refs><ref id="CVE-2020-0001" type="cve">`) {
		t.Fatalf("Marshal = %s, want the present refs element preserved", out)
	}
}

// TestGetResultsNVTRefsEmptyPresentSurvives is the presence-not-content
// counterpart to TestGetResultsNVTRefsAbsentStaysAbsent: an NVT that sends
// <refs></refs> with no <ref> children (unlike the absent case) must decode
// to a non-nil pointer, not collapse to the same nil the absent case produces.
func TestGetResultsNVTRefsEmptyPresentSurvives(t *testing.T) {
	var nvt ResultNVT
	if err := xml.Unmarshal([]byte(`<nvt><refs></refs></nvt>`), &nvt); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if nvt.Refs == nil {
		t.Fatal("Refs = nil, want a non-nil pointer for a present-but-empty <refs>")
	}
	if len(nvt.Refs.Ref) != 0 {
		t.Fatalf("Refs.Ref = %+v, want no ref entries", nvt.Refs.Ref)
	}

	out, err := xml.Marshal(&nvt)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(out), "<refs></refs>") {
		t.Fatalf("Marshal = %s, want the present-but-empty refs element preserved", out)
	}
}

// note.result and override.result are the same recursive "result inside a
// result's own notes/overrides" shape the top-level Result models with 8
// fields; the served schema itself declares this nested occurrence a plain
// leaf, so T-21 collapses it to a string instead of re-expanding the whole
// Result type recursively. Pin that collapse: attributes and any child
// elements are dropped, only the element's own character data survives.
func TestGetResultsNoteAndOverrideResultCollapseToString(t *testing.T) {
	var note ResultNote
	if err := xml.Unmarshal([]byte(`<note><result id="614a1067-8556-40ab-9b40-e2ff98a7a4f4">deleted</result></note>`), &note); err != nil {
		t.Fatalf("Unmarshal note: %v", err)
	}
	if note.Result != "deleted" {
		t.Fatalf("ResultNote.Result = %q, want %q", note.Result, "deleted")
	}

	var override ResultOverride
	if err := xml.Unmarshal([]byte(`<override><result id="614a1067-8556-40ab-9b40-e2ff98a7a4f4">deleted</result></override>`), &override); err != nil {
		t.Fatalf("Unmarshal override: %v", err)
	}
	if override.Result != "deleted" {
		t.Fatalf("ResultOverride.Result = %q, want %q", override.Result, "deleted")
	}
}
