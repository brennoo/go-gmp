package commands

import "encoding/xml"

// DeleteReportFormat represents a delete_report_format command request.
type DeleteReportFormat struct {
	XMLName        xml.Name `xml:"delete_report_format"`
	ReportFormatID string   `xml:"report_format_id,attr"`
	Ultimate       string   `xml:"ultimate,attr"`
}

// DeleteReportFormatResponse represents a delete_report_format command response.
type DeleteReportFormatResponse struct {
	XMLName    xml.Name `xml:"delete_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
