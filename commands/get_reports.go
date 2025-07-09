package commands

import "encoding/xml"

// GetReportsCommand represents a get_reports command request.
type GetReportsCommand struct {
	XMLName     xml.Name `xml:"get_reports"`
	ReportID    string   `xml:"report_id,attr,omitempty"`
	FormatID    string   `xml:"format_id,attr,omitempty"`
	Extension   string   `xml:"extension,attr,omitempty"`
	ContentType string   `xml:"content_type,attr,omitempty"`
}

// GetReportsResponse represents a get_reports command response.
type GetReportsResponse struct {
	XMLName    xml.Name        `xml:"get_reports_response"`
	Status     string          `xml:"status,attr"`
	StatusText string          `xml:"status_text,attr"`
	Reports    []ReportWrapper `xml:"report"`
}
