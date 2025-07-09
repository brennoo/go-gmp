package commands

import "encoding/xml"

// GetReportConfigs represents a get_report_configs command request.
type GetReportConfigs struct {
	XMLName        xml.Name `xml:"get_report_configs"`
	ReportConfigID string   `xml:"report_config_id,attr,omitempty"`
}

// GetReportConfigsResponse represents a get_report_configs command response.
type GetReportConfigsResponse struct {
	XMLName       xml.Name              `xml:"get_report_configs_response"`
	Status        string                `xml:"status,attr"`
	StatusText    string                `xml:"status_text,attr"`
	ReportConfigs []ReportConfigWrapper `xml:"report_config"`
}

type ReportConfigWrapper struct {
	ID               string                   `xml:"id,attr"`
	Owner            *ReportConfigOwner       `xml:"owner,omitempty"`
	Name             string                   `xml:"name"`
	Comment          string                   `xml:"comment,omitempty"`
	CreationTime     string                   `xml:"creation_time"`
	ModificationTime string                   `xml:"modification_time"`
	Writable         string                   `xml:"writable"`
	InUse            string                   `xml:"in_use"`
	Permissions      *ReportConfigPermissions `xml:"permissions,omitempty"`
	ReportFormat     *ReportConfigFormat      `xml:"report_format"`
	Params           []ReportConfigParam      `xml:"param"`
}

type ReportConfigOwner struct {
	Name string `xml:"name"`
}

type ReportConfigPermissions struct {
	Permissions []ReportConfigPermission `xml:"permission"`
}

type ReportConfigPermission struct {
	Name string `xml:"name"`
}

type ReportConfigFormat struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type ReportConfigParam struct {
	UsingDefault string   `xml:"using_default,attr,omitempty"`
	Name         string   `xml:"name"`
	Type         string   `xml:"type"`
	Value        string   `xml:"value"`
	Default      string   `xml:"default"`
	Options      []string `xml:"options>option"`
}
