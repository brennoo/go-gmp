package commands

import (
	"encoding/xml"
	"strings"
	"testing"
)

// reportIDRoundTrip unmarshals in into v, then re-marshals and checks the
// re-encoded XML against wantContains (a substring that must be present) or,
// when wantAbsent is true, must be entirely absent.
func reportIDRoundTrip(t *testing.T, in string, v any, wantAbsent bool, wantContains string) {
	t.Helper()
	if err := xml.Unmarshal([]byte(in), v); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	out, err := xml.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if wantAbsent {
		if strings.Contains(string(out), "report_id") {
			t.Fatalf("re-encoded XML fabricated <report_id>: %s", out)
		}
		return
	}
	if !strings.Contains(string(out), wantContains) {
		t.Fatalf("re-encoded XML lost the populated report_id: %s", out)
	}
}

// TestResumeStartTaskReportIDOmitEmptyFidelity is the direct behavioral proof
// finding 3 asked for: ReportID's ",omitempty" (kept via a RequiredOverrides
// entry despite the schema marking report_id required) must not fabricate an
// element on an error response, and must still round-trip a populated value
// on success.
func TestResumeStartTaskReportIDOmitEmptyFidelity(t *testing.T) {
	t.Run("resume_task error response: no report_id sent, none re-emitted", func(t *testing.T) {
		in := `<resume_task_response status="400" status_text="Task must be in Stopped or Interrupted state"></resume_task_response>`
		r := &ResumeTaskResponse{}
		reportIDRoundTrip(t, in, r, true, "")
		if r.ReportID != "" {
			t.Fatalf("ReportID = %q, want empty", r.ReportID)
		}
	})

	t.Run("resume_task success response: report_id round-trips", func(t *testing.T) {
		const id = "d4d4d4d4-1111-2222-3333-444444444444"
		in := `<resume_task_response status="202" status_text="OK"><report_id>` + id + `</report_id></resume_task_response>`
		r := &ResumeTaskResponse{}
		reportIDRoundTrip(t, in, r, false, "<report_id>"+id+"</report_id>")
		if r.ReportID != id {
			t.Fatalf("ReportID = %q, want %q", r.ReportID, id)
		}
	})

	t.Run("start_task error response: no report_id sent, none re-emitted", func(t *testing.T) {
		in := `<start_task_response status="400" status_text="Task is already active"></start_task_response>`
		r := &StartTaskResponse{}
		reportIDRoundTrip(t, in, r, true, "")
		if r.ReportID != "" {
			t.Fatalf("ReportID = %q, want empty", r.ReportID)
		}
	})

	t.Run("start_task success response: report_id round-trips", func(t *testing.T) {
		const id = "5e5e5e5e-6666-7777-8888-999999999999"
		in := `<start_task_response status="202" status_text="OK"><report_id>` + id + `</report_id></start_task_response>`
		r := &StartTaskResponse{}
		reportIDRoundTrip(t, in, r, false, "<report_id>"+id+"</report_id>")
		if r.ReportID != id {
			t.Fatalf("ReportID = %q, want %q", r.ReportID, id)
		}
	})
}
