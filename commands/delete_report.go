package commands

import "encoding/xml"

// DeleteReport represents a delete_report command request.
type DeleteReport struct {
	XMLName  xml.Name `xml:"delete_report"`
	ReportID string   `xml:"report_id,attr"`
}

// DeleteReportResponse represents a delete_report command response.
type DeleteReportResponse struct {
	XMLName    xml.Name `xml:"delete_report_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteReportResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
