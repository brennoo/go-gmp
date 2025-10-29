package commands

import "encoding/xml"

// CreateReport represents a create_report command request.
type CreateReport struct {
	XMLName  xml.Name    `xml:"create_report"`
	Report   *Report     `xml:"report"`
	Task     *ReportTask `xml:"task,omitempty"`
	InAssets *bool       `xml:"in_assets,omitempty"`
}

// CreateReportResponse represents a create_report command response.
type CreateReportResponse struct {
	XMLName    xml.Name `xml:"create_report_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateReportResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
