package commands

import "encoding/xml"

// ModifyReportConfig represents a modify_report_config command request.
type ModifyReportConfig struct {
	XMLName xml.Name                  `xml:"modify_report_config"`
	Name    string                    `xml:"name,omitempty"`
	Comment string                    `xml:"comment,omitempty"`
	Params  []ModifyReportConfigParam `xml:"param,omitempty"`
}

type ModifyReportConfigParam struct {
	Name       string `xml:"name"`
	Value      string `xml:"value"`
	UseDefault string `xml:"use_default,attr,omitempty"`
}

// ModifyReportConfigResponse represents a modify_report_config command response.
type ModifyReportConfigResponse struct {
	XMLName    xml.Name `xml:"modify_report_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
