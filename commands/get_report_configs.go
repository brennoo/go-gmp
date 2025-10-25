package commands

import (
	"encoding/xml"
	"time"
)

// GetReportConfigs represents a get_report_configs command request.
type GetReportConfigs struct {
	XMLName        xml.Name `xml:"get_report_configs"`
	ReportConfigID string   `xml:"report_config_id,attr,omitempty"`
}

// GetReportConfigsResponse represents a get_report_configs command response.
type GetReportConfigsResponse struct {
	XMLName       xml.Name                `xml:"get_report_configs_response"`
	Status        string                  `xml:"status,attr"`
	StatusText    string                  `xml:"status_text,attr"`
	ReportConfigs []ReportConfig          `xml:"report_config"`
	Filters       *ReportConfigFilters    `xml:"filters,omitempty"`
	Sort          *ReportConfigSort       `xml:"sort,omitempty"`
	Pagination    *ReportConfigPagination `xml:"report_configs,omitempty"`
	Count         *ReportConfigCount      `xml:"report_config_count,omitempty"`
}

// ReportConfig represents a report_config element in the response.
type ReportConfig struct {
	ID               string                    `xml:"id,attr"`
	Owner            *ReportConfigOwner        `xml:"owner,omitempty"`
	Name             string                    `xml:"name"`
	Comment          string                    `xml:"comment,omitempty"`
	CreationTime     time.Time                 `xml:"creation_time"`
	ModificationTime time.Time                 `xml:"modification_time"`
	Writable         bool                      `xml:"writable"`
	InUse            bool                      `xml:"in_use"`
	Permissions      *ReportConfigPermissions  `xml:"permissions,omitempty"`
	ReportFormat     *ReportConfigReportFormat `xml:"report_format"`
	Params           []ReportConfigParam       `xml:"param"`
}

// ReportConfigOwner represents the owner of a report config.
type ReportConfigOwner struct {
	Name string `xml:"name"`
}

// ReportConfigPermissions represents permissions for a report config.
type ReportConfigPermissions struct {
	Permissions []ReportConfigPermission `xml:"permission"`
}

// ReportConfigPermission represents a single permission for a report config.
type ReportConfigPermission struct {
	Name string `xml:"name"`
}

// ReportConfigReportFormat represents the report_format element in a report config.
type ReportConfigReportFormat struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// ReportConfigParam represents a param element in a report config.
type ReportConfigParam struct {
	UsingDefault bool     `xml:"using_default,attr,omitempty"`
	Name         string   `xml:"name"`
	Type         string   `xml:"type"`
	Value        string   `xml:"value"`
	Default      string   `xml:"default"`
	Options      []string `xml:"options>option"`
}

// ReportConfigFilters represents filters for the get_report_configs response.
type ReportConfigFilters struct {
	ID       string `xml:"id,attr,omitempty"`
	Term     string `xml:"term,omitempty"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

// ReportConfigSort represents sorting information for the get_report_configs response.
type ReportConfigSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// ReportConfigPagination represents pagination information for the get_report_configs response.
type ReportConfigPagination struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// ReportConfigCount represents the report_config_count element in the get_report_configs response.
type ReportConfigCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
