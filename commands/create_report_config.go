package commands

import "encoding/xml"

// CreateReportConfig is a create_report_config command request.
type CreateReportConfig struct {
	XMLName      xml.Name                 `xml:"create_report_config"`
	Name         string                   `xml:"name"`
	Comment      string                   `xml:"comment,omitempty"`
	Copy         string                   `xml:"copy,omitempty"`
	ReportFormat ReportConfigReportFormat `xml:"report_format"`
	Params       []ReportConfigParam      `xml:"param,omitempty"`
}

// CreateReportConfigResponse is a create_report_config command response.
type CreateReportConfigResponse struct {
	XMLName    xml.Name `xml:"create_report_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
