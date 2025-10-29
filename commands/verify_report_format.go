package commands

import "encoding/xml"

// VerifyReportFormat represents a verify_report_format command request.
type VerifyReportFormat struct {
	XMLName        xml.Name `xml:"verify_report_format"`
	ReportFormatID string   `xml:"report_format_id,attr"`
}

// VerifyReportFormatResponse represents a verify_report_format command response.
type VerifyReportFormatResponse struct {
	XMLName    xml.Name `xml:"verify_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *VerifyReportFormatResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
