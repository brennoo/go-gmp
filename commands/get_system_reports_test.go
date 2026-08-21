package commands

import (
	"encoding/xml"
	"strings"
	"testing"
)

// Present start_time/end_time, including the year-1 zero value, round-trip verbatim.
func TestGetSystemReportsReportTimestampsPresent(t *testing.T) {
	in := `<get_system_reports_response status="200" status_text="OK">` +
		`<system_report><name>fallback</name><title>Fallback Report</title>` +
		`<report format="txt" start_time="0001-01-01T00:00:00Z" end_time="0001-01-01T00:00:00Z" duration="">text</report>` +
		`</system_report></get_system_reports_response>`
	var r GetSystemReportsResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if len(r.SystemReport) != 1 || r.SystemReport[0].Report == nil {
		t.Fatalf("SystemReport = %+v, want one entry with a Report", r.SystemReport)
	}
	rep := r.SystemReport[0].Report
	if rep.StartTime == nil || !rep.StartTime.IsZero() {
		t.Fatalf("StartTime = %v, want a non-nil pointer to the zero time", rep.StartTime)
	}
	if rep.EndTime == nil || !rep.EndTime.IsZero() {
		t.Fatalf("EndTime = %v, want a non-nil pointer to the zero time", rep.EndTime)
	}

	out, err := xml.Marshal(&r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(out); !strings.Contains(got, `start_time="0001-01-01T00:00:00Z"`) || !strings.Contains(got, `end_time="0001-01-01T00:00:00Z"`) {
		t.Fatalf("Marshal = %s, want both timestamps preserved verbatim", got)
	}
}

// Absent start_time/end_time must stay nil, not fabricate a zero timestamp on marshal.
func TestGetSystemReportsReportTimestampsAbsent(t *testing.T) {
	in := `<get_system_reports_response status="200" status_text="OK">` +
		`<system_report><name>proc</name><title>Processes</title>` +
		`<report format="png" duration="86400">aW1hZ2U=</report>` +
		`</system_report></get_system_reports_response>`
	var r GetSystemReportsResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	rep := r.SystemReport[0].Report
	if rep.StartTime != nil {
		t.Fatalf("StartTime = %v, want nil since the source never sent it", rep.StartTime)
	}
	if rep.EndTime != nil {
		t.Fatalf("EndTime = %v, want nil since the source never sent it", rep.EndTime)
	}

	out, err := xml.Marshal(&r)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(out); strings.Contains(got, "start_time") || strings.Contains(got, "end_time") {
		t.Fatalf("Marshal = %s, want no fabricated start_time/end_time", got)
	}
}

// A request that never sets StartTime/EndTime must not fabricate a zero
// timestamp on marshal (this is the request-side counterpart of
// TestGetSystemReportsReportTimestampsAbsent above, catching the same
// non-pointer-GMPTime-plus-omitempty bug on the request's own attributes).
func TestGetSystemReportsRequestTimestampsAbsent(t *testing.T) {
	req := GetSystemReports{Name: "load"}

	out, err := xml.Marshal(&req)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(out); strings.Contains(got, "start_time") || strings.Contains(got, "end_time") {
		t.Fatalf("Marshal = %s, want no fabricated start_time/end_time", got)
	}
}
