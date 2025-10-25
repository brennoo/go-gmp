package commands

import (
	"encoding/xml"
	"time"
)

// GetNVTs represents a get_nvts command request.
type GetNVTs struct {
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

// GetNVTsResponse represents a get_nvts command response.
type GetNVTsResponse struct {
	XMLName    xml.Name `xml:"get_nvts_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	NVTs       []NVT    `xml:"nvt,omitempty"`
}

// NVT represents a Network Vulnerability Test in the response.
type NVT struct {
	OID              string          `xml:"oid,attr"`
	Name             string          `xml:"name,omitempty"`
	CreationTime     time.Time       `xml:"creation_time,omitempty"`
	ModificationTime time.Time       `xml:"modification_time,omitempty"`
	Category         int             `xml:"category,omitempty"`
	Family           string          `xml:"family,omitempty"`
	CvssBase         int             `xml:"cvss_base,omitempty"`
	Severities       *NVTSeverities  `xml:"severities,omitempty"`
	Refs             string          `xml:"refs,omitempty"`
	Tags             string          `xml:"tags,omitempty"`
	PreferenceCount  int             `xml:"preference_count,omitempty"`
	Timeout          int             `xml:"timeout,omitempty"`
	DefaultTimeout   int             `xml:"default_timeout,omitempty"`
	Solution         string          `xml:"solution,omitempty"`
	Epss             string          `xml:"epss,omitempty"`
	Preferences      *NVTPreferences `xml:"preferences,omitempty"`
	UserTags         string          `xml:"user_tags,omitempty"`
	Qod              string          `xml:"qod,omitempty"`
	Description      string          `xml:"description,omitempty"`
}

// NVTSeverities represents the severities element in an NVT.
type NVTSeverities struct {
	Score      float64       `xml:"score,attr,omitempty"`
	Severities []NVTSeverity `xml:"severity,omitempty"`
}

// NVTSeverity represents a severity element in an NVT.
type NVTSeverity struct {
	Type   string  `xml:"type,attr,omitempty"`
	Origin string  `xml:"origin,omitempty"`
	Date   string  `xml:"date,omitempty"`
	Score  float64 `xml:"score,omitempty"`
	Value  string  `xml:"value,omitempty"`
}

// NVTPreferences represents the preferences element in an NVT.
type NVTPreferences struct {
	Timeout        string          `xml:"timeout,omitempty"`
	DefaultTimeout string          `xml:"default_timeout,omitempty"`
	Preferences    []NVTPreference `xml:"preference,omitempty"`
}

// NVTPreference represents a preference element in an NVT.
type NVTPreference struct {
	NVT     *NVTPreferenceNVT `xml:"nvt,omitempty"`
	ID      string            `xml:"id,omitempty"`
	Name    string            `xml:"name,omitempty"`
	HRName  string            `xml:"hr_name,omitempty"`
	Type    string            `xml:"type,omitempty"`
	Value   string            `xml:"value,omitempty"`
	Alt     string            `xml:"alt,omitempty"`
	Default string            `xml:"default,omitempty"`
}

// NVTPreferenceNVT represents an NVT reference in a preference.
type NVTPreferenceNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name,omitempty"`
}
