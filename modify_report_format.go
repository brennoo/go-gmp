package gmp

import "encoding/xml"

// ModifyReportFormatCommand represents a modify_report_format command request.
type ModifyReportFormatCommand struct {
	XMLName        xml.Name                 `xml:"modify_report_format"`
	ReportFormatID string                   `xml:"report_format_id,attr"`
	Active         *bool                    `xml:"active,omitempty"`
	Name           string                   `xml:"name,omitempty"`
	Summary        string                   `xml:"summary,omitempty"`
	Param          *ModifyReportFormatParam `xml:"param,omitempty"`
}

// ModifyReportFormatParam represents the <param> element for modify_report_format.
type ModifyReportFormatParam struct {
	Name  string `xml:"name"`
	Value string `xml:"value,omitempty"`
}

// ModifyReportFormatResponse represents a modify_report_format command response.
type ModifyReportFormatResponse struct {
	XMLName    xml.Name `xml:"modify_report_format_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
