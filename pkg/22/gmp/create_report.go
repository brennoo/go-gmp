package gmp

import "encoding/xml"

// CreateReportCommand represents a create_report command request.
type CreateReportCommand struct {
	XMLName  xml.Name          `xml:"create_report"`
	Report   *ReportWrapper    `xml:"report"`
	Task     *CreateReportTask `xml:"task,omitempty"`
	InAssets *bool             `xml:"in_assets,omitempty"`
}

// ReportWrapper represents the <report> element for create_report.
type ReportWrapper struct {
	ID          string `xml:"id,attr,omitempty"`
	FormatID    string `xml:"format_id,attr,omitempty"`
	Extension   string `xml:"extension,attr,omitempty"`
	ContentType string `xml:"content_type,attr,omitempty"`
	// The actual report content is complex and can be nested; for now, use RawXML for flexibility.
	RawXML []byte `xml:",innerxml"`
}

// CreateReportTask represents the <task> element for create_report.
type CreateReportTask struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateReportResponse represents a create_report command response.
type CreateReportResponse struct {
	XMLName    xml.Name `xml:"create_report_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
