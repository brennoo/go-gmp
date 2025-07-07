package gmp

import "encoding/xml"

// CreateReportConfigCommand represents a create_report_config command request.
type CreateReportConfigCommand struct {
	XMLName      xml.Name                  `xml:"create_report_config"`
	Name         string                    `xml:"name"`
	Comment      string                    `xml:"comment,omitempty"`
	Copy         string                    `xml:"copy,omitempty"`
	ReportFormat CreateReportConfigFormat  `xml:"report_format"`
	Params       []CreateReportConfigParam `xml:"param,omitempty"`
}

type CreateReportConfigFormat struct {
	ID string `xml:"id,attr"`
}

type CreateReportConfigParam struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// CreateReportConfigResponse represents a create_report_config command response.
type CreateReportConfigResponse struct {
	XMLName    xml.Name `xml:"create_report_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
