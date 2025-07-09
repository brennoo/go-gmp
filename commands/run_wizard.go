package commands

import (
	"encoding/xml"
)

// RunWizardCommand represents a run_wizard command request.
type RunWizardCommand struct {
	XMLName  xml.Name      `xml:"run_wizard"`
	Mode     string        `xml:"mode,omitempty"`
	Name     string        `xml:"name"`
	Params   *WizardParams `xml:"params,omitempty"`
	ReadOnly string        `xml:"read_only,attr,omitempty"`
}

type WizardParams struct {
	Params []WizardParam `xml:"param"`
}

type WizardParam struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// RunWizardResponse represents a run_wizard command response.
type RunWizardResponse struct {
	XMLName    xml.Name        `xml:"run_wizard_response"`
	Status     string          `xml:"status,attr"`
	StatusText string          `xml:"status_text,attr"`
	Response   *WizardResponse `xml:"response,omitempty"`
}

type WizardResponse struct {
	StartTaskResponse *StartTaskResponse `xml:"start_task_response,omitempty"`
	// TODO: check other response types
}
