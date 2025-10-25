package commands

import (
	"encoding/xml"
	"time"
)

// GetVulns represents a get_vulns command request.
type GetVulns struct {
	XMLName xml.Name `xml:"get_vulns"`
	VulnID  string   `xml:"vuln_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
	FiltID  string   `xml:"filt_id,attr,omitempty"`
	Trash   bool     `xml:"trash,attr,omitempty"`
}

// GetVulnsResponse represents a get_vulns command response.
type GetVulnsResponse struct {
	XMLName    xml.Name      `xml:"get_vulns_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Vulns      []Vuln        `xml:"vuln,omitempty"`
	Filters    *VulnsFilters `xml:"filters,omitempty"`
	Sort       *VulnsSort    `xml:"sort,omitempty"`
	VulnsMeta  *VulnsMeta    `xml:"vulns,omitempty"`
	VulnCount  *VulnCount    `xml:"vuln_count,omitempty"`
	Truncated  string        `xml:"truncated,omitempty"`
}

// Vuln represents a <vuln> element in the get_vulns response.
type Vuln struct {
	XMLName          xml.Name     `xml:"vuln"`
	ID               string       `xml:"id,attr"`
	Name             string       `xml:"name"`
	Type             string       `xml:"type"`
	CreationTime     time.Time    `xml:"creation_time"`
	ModificationTime time.Time    `xml:"modification_time"`
	Severity         float64      `xml:"severity"`
	QoD              int          `xml:"qod"`
	Results          *VulnResults `xml:"results,omitempty"`
	Hosts            *VulnHosts   `xml:"hosts,omitempty"`
}

// VulnResults represents the <results> element in a vuln.
type VulnResults struct {
	Count  int       `xml:"count"`
	Oldest time.Time `xml:"oldest"`
	Newest time.Time `xml:"newest"`
}

// VulnHosts represents the <hosts> element in a vuln.
type VulnHosts struct {
	Count int `xml:"count"`
}

// VulnsFilters represents the <filters> element in the response.
type VulnsFilters struct {
	XMLName  xml.Name       `xml:"filters"`
	ID       string         `xml:"id,attr"`
	Term     string         `xml:"term"`
	Name     string         `xml:"name,omitempty"`
	Keywords *VulnsKeywords `xml:"keywords,omitempty"`
}

// VulnsKeywords represents the <keywords> element in filters.
type VulnsKeywords struct {
	Keyword []VulnsKeyword `xml:"keyword"`
}

// VulnsKeyword represents a <keyword> element in filters keywords.
type VulnsKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// VulnsSort represents the <sort> element in the response.
type VulnsSort struct {
	Field *VulnsSortField `xml:"field,omitempty"`
}

// VulnsSortField represents the <field> element in sort.
type VulnsSortField struct {
	Order string `xml:"order"`
}

// VulnsMeta represents the <vulns> meta element in the response.
type VulnsMeta struct {
	XMLName xml.Name `xml:"vulns"`
	Start   int      `xml:"start,attr"`
	Max     int      `xml:"max,attr"`
}

// VulnCount represents the <vuln_count> element in the response.
type VulnCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
