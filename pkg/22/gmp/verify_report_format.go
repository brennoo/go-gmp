package gmp

import "encoding/xml"

// VerifyReportFormatCommand represents a verify_report_format command request.
type VerifyReportFormatCommand struct {
	XMLName        xml.Name `xml:"verify_report_format"`
	ReportFormatID string   `xml:"report_format_id,attr"`
}

// VerifyReportFormatResponse represents a verify_report_format command response.
type VerifyReportFormatResponse struct {
	XMLName    xml.Name `xml:"verify_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
