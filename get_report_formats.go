package gmp

import "encoding/xml"

// GetReportFormatsCommand represents a get_report_formats command request.
type GetReportFormatsCommand struct {
	XMLName        xml.Name `xml:"get_report_formats"`
	ReportFormatID string   `xml:"report_format_id,attr,omitempty"`
	Details        *int     `xml:"details,attr,omitempty"`
}

// ReportFormatWrapper represents the <report_format> element in the response.
type ReportFormatWrapper struct {
	ID               string                 `xml:"id,attr"`
	Name             string                 `xml:"name"`
	Comment          string                 `xml:"comment,omitempty"`
	CreationTime     string                 `xml:"creation_time"`
	ModificationTime string                 `xml:"modification_time"`
	Writable         string                 `xml:"writable"`
	InUse            string                 `xml:"in_use"`
	Extension        string                 `xml:"extension"`
	ContentType      string                 `xml:"content_type"`
	Summary          string                 `xml:"summary"`
	Description      string                 `xml:"description"`
	Trust            *ReportFormatTrust     `xml:"trust,omitempty"`
	Active           string                 `xml:"active"`
	File             []ReportFormatFile     `xml:"file"`
	Signature        *ReportFormatSignature `xml:"signature,omitempty"`
}

type ReportFormatTrust struct {
	Value string `xml:",chardata"`
	Time  string `xml:"time"`
}

type ReportFormatFile struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type ReportFormatSignature struct {
	Value string `xml:",chardata"`
}

// GetReportFormatsResponse represents a get_report_formats command response.
type GetReportFormatsResponse struct {
	XMLName       xml.Name              `xml:"get_report_formats_response"`
	Status        string                `xml:"status,attr"`
	StatusText    string                `xml:"status_text,attr"`
	ReportFormats []ReportFormatWrapper `xml:"report_format"`
}
