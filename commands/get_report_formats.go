package commands

import "encoding/xml"

// GetReportFormats represents a get_report_formats command request.
type GetReportFormats struct {
	XMLName        xml.Name `xml:"get_report_formats"`
	ReportFormatID string   `xml:"report_format_id,attr,omitempty"`
	Details        *int     `xml:"details,attr,omitempty"`
}

// GetReportFormatsResponse represents a get_report_formats command response.
type GetReportFormatsResponse struct {
	XMLName       xml.Name                 `xml:"get_report_formats_response"`
	Status        string                   `xml:"status,attr"`
	StatusText    string                   `xml:"status_text,attr"`
	ReportFormats []ReportFormat           `xml:"report_format"`
	Filters       *ReportFormatsFilters    `xml:"filters,omitempty"`
	Sort          *ReportFormatsSort       `xml:"sort,omitempty"`
	Pagination    *ReportFormatsPagination `xml:"report_formats,omitempty"`
	Count         *ReportFormatsCount      `xml:"report_format_count,omitempty"`
}

// ReportFormat represents a report_format element in the response.
type ReportFormat struct {
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
	Files            []ReportFormatFile     `xml:"file"`
	Signature        *ReportFormatSignature `xml:"signature,omitempty"`
	Predefined       string                 `xml:"predefined,omitempty"`
	Configurable     string                 `xml:"configurable,omitempty"`
	ReportType       string                 `xml:"report_type,omitempty"`
	Deprecated       string                 `xml:"deprecated,omitempty"`
}

// ReportFormatTrust represents the trust element in a report format.
type ReportFormatTrust struct {
	Value string `xml:",chardata"`
	Time  string `xml:"time"`
}

// ReportFormatFile represents a file element in a report format.
type ReportFormatFile struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// ReportFormatSignature represents the signature element in a report format.
type ReportFormatSignature struct {
	Value string `xml:",chardata"`
}

// ReportFormatsFilters represents filters for the get_report_formats response.
type ReportFormatsFilters struct {
	ID       string `xml:"id,attr,omitempty"`
	Term     string `xml:"term,omitempty"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

// ReportFormatsSort represents sorting information for the get_report_formats response.
type ReportFormatsSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// ReportFormatsPagination represents pagination information for the get_report_formats response.
type ReportFormatsPagination struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// ReportFormatsCount represents the report_format_count element in the get_report_formats response.
type ReportFormatsCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
