package commands

import "encoding/xml"

// CreateReportFormat represents a create_report_format command request.
type CreateReportFormat struct {
	XMLName xml.Name `xml:"create_report_format"`
	Name    string   `xml:"name,omitempty"`
	Copy    string   `xml:"copy,omitempty"`
}

// CreateReportFormatResponse represents a create_report_format command response.
type CreateReportFormatResponse struct {
	XMLName    xml.Name `xml:"create_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateReportFormatResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
