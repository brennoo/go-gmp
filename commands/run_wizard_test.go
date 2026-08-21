package commands

import (
	"encoding/xml"
	"strings"
	"testing"
)

// A wizard-embeddable response gvmd happens to enumerate in WizardResponse's
// typed fields decodes into that field.
func TestWizardResponseKnownTypeDecodesTyped(t *testing.T) {
	in := `<run_wizard_response status="202" status_text="OK, request submitted">` +
		`<response><start_task_response status="202" status_text="OK, request submitted">` +
		`<report_id>a06d21f7-8e2f-4d7f-bceb-1df852d8b37d</report_id></start_task_response></response>` +
		`</run_wizard_response>`

	var r RunWizardResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r.Response == nil || r.Response.StartTaskResponse == nil {
		t.Fatalf("Response.StartTaskResponse = nil, want decoded")
	}
	if got := r.Response.StartTaskResponse.ReportID; got != "a06d21f7-8e2f-4d7f-bceb-1df852d8b37d" {
		t.Fatalf("ReportID = %q, want the report id preserved", got)
	}
}

// WizardResponse's ~35 typed fields are hand-enumerated from gvmd's known
// wizard scripts, not schema-derived (the schema declares <response> as an
// empty pattern with no fixed content model), so gvmd can embed a response
// type this list doesn't cover. CreateFilterResponse is one such gap.
// InnerXML must still preserve it verbatim so no wizard-embedded content is
// ever silently lost, even when no typed field claims it.
func TestWizardResponseUnknownTypePreservedViaInnerXML(t *testing.T) {
	in := `<run_wizard_response status="202" status_text="OK, request submitted">` +
		`<response><create_filter_response status="201" status_text="OK, resource created" id="abc-123"/></response>` +
		`</run_wizard_response>`

	var r RunWizardResponse
	if err := xml.Unmarshal([]byte(in), &r); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if r.Response == nil {
		t.Fatalf("Response = nil, want decoded")
	}
	if !strings.Contains(r.Response.InnerXML, `<create_filter_response`) || !strings.Contains(r.Response.InnerXML, `id="abc-123"`) {
		t.Fatalf("InnerXML = %q, want the unenumerated response type preserved verbatim", r.Response.InnerXML)
	}
}
