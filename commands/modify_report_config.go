package commands

import "encoding/xml"

// ModifyReportConfig represents a modify_report_config command request.
type ModifyReportConfig struct {
	XMLName xml.Name            `xml:"modify_report_config"`
	Name    string              `xml:"name,omitempty"`
	Comment string              `xml:"comment,omitempty"`
	Params  []ReportConfigParam `xml:"param,omitempty"`
}

// ModifyReportConfigResponse represents a modify_report_config command response.
type ModifyReportConfigResponse struct {
	XMLName    xml.Name `xml:"modify_report_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyReportConfigResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
