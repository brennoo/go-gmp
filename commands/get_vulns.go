package commands

import "encoding/xml"

// GetVulns represents a get_vulns command request.
type GetVulns struct {
	XMLName string `xml:"get_vulns"`
	VulnID  string `xml:"vuln_id,attr,omitempty"`
	Filter  string `xml:"filter,attr,omitempty"`
	FiltID  string `xml:"filt_id,attr,omitempty"`
	Trash   bool   `xml:"trash,attr,omitempty"`
}

// GetVulnsResponse represents a get_vulns command response.
type GetVulnsResponse struct {
	XMLName    xml.Name           `xml:"get_vulns_response"`
	Status     string             `xml:"status,attr"`
	StatusText string             `xml:"status_text,attr"`
	Vuln       []*GetVulnsVuln    `xml:"vuln,omitempty"`
	Filters    *GetVulnsFilters   `xml:"filters,omitempty"`
	Sort       *GetVulnsSort      `xml:"sort,omitempty"`
	Vulns      *GetVulnsVulns     `xml:"vulns,omitempty"`
	VulnCount  *GetVulnsVulnCount `xml:"vuln_count,omitempty"`
}

type GetVulnsVuln struct {
	XMLName          xml.Name             `xml:"vuln"`
	ID               string               `xml:"id,attr,omitempty"`
	Name             string               `xml:"name,omitempty"`
	Type             string               `xml:"type,omitempty"`
	CreationTime     string               `xml:"creation_time,omitempty"`
	ModificationTime string               `xml:"modification_time,omitempty"`
	Severity         string               `xml:"severity,omitempty"`
	QoD              int                  `xml:"qod,omitempty"`
	Results          *GetVulnsVulnResults `xml:"results,omitempty"`
	Hosts            *GetVulnsVulnHosts   `xml:"hosts,omitempty"`
}

type GetVulnsVulnResults struct {
	Count  int    `xml:"count,omitempty"`
	Oldest string `xml:"oldest,omitempty"`
	Newest string `xml:"newest,omitempty"`
}

type GetVulnsVulnHosts struct {
	Count int `xml:"count,omitempty"`
}

type GetVulnsFilters struct {
	XMLName  xml.Name                 `xml:"filters"`
	ID       string                   `xml:"id,attr,omitempty"`
	Term     string                   `xml:"term,omitempty"`
	Name     string                   `xml:"name,omitempty"`
	Keywords *GetVulnsFiltersKeywords `xml:"keywords,omitempty"`
}

type GetVulnsFiltersKeywords struct {
	Keyword []*GetVulnsFiltersKeyword `xml:"keyword,omitempty"`
}

type GetVulnsFiltersKeyword struct {
	Column   string `xml:"column,omitempty"`
	Relation string `xml:"relation,omitempty"`
	Value    string `xml:"value,omitempty"`
}

type GetVulnsSort struct {
	Field *GetVulnsSortField `xml:"field,omitempty"`
}

type GetVulnsSortField struct {
	Order string `xml:"order,omitempty"`
}

type GetVulnsVulns struct {
	XMLName xml.Name `xml:"vulns"`
	Start   int      `xml:"start,attr,omitempty"`
	Max     int      `xml:"max,attr,omitempty"`
}

type GetVulnsVulnCount struct {
	Filtered int `xml:"filtered,omitempty"`
	Page     int `xml:"page,omitempty"`
}
