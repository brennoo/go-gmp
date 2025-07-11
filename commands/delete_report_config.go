package commands

import "encoding/xml"

// DeleteReportConfig represents a delete_report_config command request.
type DeleteReportConfig struct {
	XMLName        xml.Name `xml:"delete_report_config"`
	ReportConfigID string   `xml:"report_config_id,attr"`
	Ultimate       string   `xml:"ultimate,attr"`
}

// DeleteReportConfigResponse represents a delete_report_config command response.
type DeleteReportConfigResponse struct {
	XMLName    xml.Name `xml:"delete_report_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
