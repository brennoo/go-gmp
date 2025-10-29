package commands

import "encoding/xml"

// ModifyReportFormat represents a modify_report_format command request.
type ModifyReportFormat struct {
	XMLName        xml.Name           `xml:"modify_report_format"`
	ReportFormatID string             `xml:"report_format_id,attr"`
	Active         *bool              `xml:"active,omitempty"`
	Name           string             `xml:"name,omitempty"`
	Summary        string             `xml:"summary,omitempty"`
	Param          *ReportFormatParam `xml:"param,omitempty"`
}

// ModifyReportFormatResponse represents a modify_report_format command response.
type ModifyReportFormatResponse struct {
	XMLName    xml.Name `xml:"modify_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyReportFormatResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
