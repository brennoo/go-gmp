package commands

import "encoding/xml"

// GetNvtsCommand represents a get_nvts command request.
type GetNvtsCommand struct {
	XMLName             xml.Name `xml:"get_nvts"`
	NvtOID              string   `xml:"nvt_oid,attr,omitempty"`
	Details             string   `xml:"details,attr,omitempty"`
	Lean                string   `xml:"lean,attr,omitempty"`
	Preferences         string   `xml:"preferences,attr,omitempty"`
	PreferenceCount     string   `xml:"preference_count,attr,omitempty"`
	SkipCertRefs        string   `xml:"skip_cert_refs,attr,omitempty"`
	SkipTags            string   `xml:"skip_tags,attr,omitempty"`
	Timeout             string   `xml:"timeout,attr,omitempty"`
	ConfigID            string   `xml:"config_id,attr,omitempty"`
	PreferencesConfigID string   `xml:"preferences_config_id,attr,omitempty"`
	Family              string   `xml:"family,attr,omitempty"`
	SortOrder           string   `xml:"sort_order,attr,omitempty"`
	SortField           string   `xml:"sort_field,attr,omitempty"`
}

// GetNvtsResponse represents a get_nvts command response.
type GetNvtsResponse struct {
	XMLName    xml.Name `xml:"get_nvts_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Nvts       []NVT    `xml:"nvt,omitempty"`
}

// NVT represents a Network Vulnerability Test.
type NVT struct {
	OID              string `xml:"oid,attr"`
	Name             string `xml:"name,omitempty"`
	UserTags         string `xml:"user_tags,omitempty"`
	CreationTime     string `xml:"creation_time,omitempty"`
	ModificationTime string `xml:"modification_time,omitempty"`
	Category         string `xml:"category,omitempty"`
	Family           string `xml:"family,omitempty"`
	CvssBase         string `xml:"cvss_base,omitempty"`
	Severities       string `xml:"severities,omitempty"`
	Qod              string `xml:"qod,omitempty"`
	Refs             string `xml:"refs,omitempty"`
	Tags             string `xml:"tags,omitempty"`
	PreferenceCount  string `xml:"preference_count,omitempty"`
	Timeout          string `xml:"timeout,omitempty"`
	DefaultTimeout   string `xml:"default_timeout,omitempty"`
	Solution         string `xml:"solution,omitempty"`
	Epss             string `xml:"epss,omitempty"`
	Preferences      string `xml:"preferences,omitempty"`
}
