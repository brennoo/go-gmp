package commands

import (
	"encoding/xml"
	"time"
)

// GetOverrides represents a get_overrides command request.
type GetOverrides struct {
	XMLName    xml.Name `xml:"get_overrides"`
	OverrideID string   `xml:"override_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
	FiltID     string   `xml:"filt_id,attr,omitempty"`
	Trash      string   `xml:"trash,attr,omitempty"`
	Details    string   `xml:"details,attr,omitempty"`
	Result     string   `xml:"result,attr,omitempty"`
}

// GetOverridesResponse represents a get_overrides command response.
type GetOverridesResponse struct {
	XMLName    xml.Name            `xml:"get_overrides_response"`
	Status     string              `xml:"status,attr"`
	StatusText string              `xml:"status_text,attr"`
	Overrides  []Override          `xml:"override,omitempty"`
	Filters    *OverrideFilters    `xml:"filters,omitempty"`
	Sort       *OverrideSort       `xml:"sort,omitempty"`
	Pagination *OverridePagination `xml:"overrides,omitempty"`
	Count      *OverrideCount      `xml:"override_count,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetOverridesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Override represents an override in the response.
type Override struct {
	ID               string               `xml:"id,attr"`
	Permissions      *OverridePermissions `xml:"permissions,omitempty"`
	Owner            *OverrideOwner       `xml:"owner,omitempty"`
	NVT              *OverrideNVT         `xml:"nvt,omitempty"`
	CreationTime     time.Time            `xml:"creation_time,omitempty"`
	ModificationTime time.Time            `xml:"modification_time,omitempty"`
	Writable         bool                 `xml:"writable,omitempty"`
	InUse            bool                 `xml:"in_use,omitempty"`
	Active           bool                 `xml:"active,omitempty"`
	Text             *OverrideText        `xml:"text,omitempty"`
	Threat           string               `xml:"threat,omitempty"`
	Severity         float64              `xml:"severity,omitempty"`
	NewThreat        string               `xml:"new_threat,omitempty"`
	NewSeverity      float64              `xml:"new_severity,omitempty"`
	Orphan           bool                 `xml:"orphan,omitempty"`
	UserTags         *OverrideUserTags    `xml:"user_tags,omitempty"`
	Hosts            string               `xml:"hosts,omitempty"`
	Port             string               `xml:"port,omitempty"`
	Task             *OverrideTask        `xml:"task,omitempty"`
	EndTime          string               `xml:"end_time,omitempty"`
	Result           *OverrideResult      `xml:"result,omitempty"`
}

// OverridePermissions represents permissions information in an override.
type OverridePermissions struct {
	Permission []OverridePermission `xml:"permission,omitempty"`
}

// OverridePermission represents a single permission in an override.
type OverridePermission struct {
	Name string `xml:"name"`
}

// OverrideOwner represents owner information in an override.
type OverrideOwner struct {
	Name string `xml:"name"`
}

// OverrideNVT represents NVT information in an override.
type OverrideNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name,omitempty"`
	Type string `xml:"type,omitempty"`
}

// OverrideText represents text information in an override.
type OverrideText struct {
	Text    string `xml:",chardata"`
	Excerpt string `xml:"excerpt,attr,omitempty"`
}

// OverrideUserTags represents user tags information in an override.
type OverrideUserTags struct {
	Count int               `xml:"count,omitempty"`
	Tag   []OverrideUserTag `xml:"tag,omitempty"`
}

// OverrideUserTag represents a user tag in an override.
type OverrideUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// OverrideTask represents task information in an override.
type OverrideTask struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name,omitempty"`
	Trash string `xml:"trash,omitempty"`
}

// OverrideResult represents result information in an override.
type OverrideResult struct {
	ID          string              `xml:"id,attr"`
	Host        *OverrideResultHost `xml:"host,omitempty"`
	Port        string              `xml:"port,omitempty"`
	NVT         *OverrideResultNVT  `xml:"nvt,omitempty"`
	Threat      string              `xml:"threat,omitempty"`
	Severity    string              `xml:"severity,omitempty"`
	QOD         *OverrideResultQOD  `xml:"qod,omitempty"`
	Description string              `xml:"description,omitempty"`
}

// OverrideResultHost represents host information in an override result.
type OverrideResultHost struct {
	Value string                   `xml:",chardata"`
	Asset *OverrideResultHostAsset `xml:"asset,omitempty"`
}

// OverrideResultHostAsset represents asset information in an override result host.
type OverrideResultHostAsset struct {
	AssetID string `xml:"asset_id,attr"`
}

// OverrideResultNVT represents NVT information in an override result.
type OverrideResultNVT struct {
	OID      string `xml:"oid,attr"`
	Name     string `xml:"name,omitempty"`
	Type     string `xml:"type,omitempty"`
	CVSSBase string `xml:"cvss_base,omitempty"`
	CVE      string `xml:"cve,omitempty"`
	BID      int    `xml:"bid,omitempty"`
}

// OverrideResultQOD represents QOD information in an override result.
type OverrideResultQOD struct {
	Value int    `xml:"value,omitempty"`
	Type  string `xml:"type,omitempty"`
}

// OverrideFilters represents filter information in the response.
type OverrideFilters struct {
	ID       string                   `xml:"id,attr"`
	Term     string                   `xml:"term"`
	Name     string                   `xml:"name,omitempty"`
	Keywords []OverrideFiltersKeyword `xml:"keywords>keyword,omitempty"`
}

// OverrideFiltersKeyword represents a keyword in override filters.
type OverrideFiltersKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// OverrideSort represents sort information in the response.
type OverrideSort struct {
	Field OverrideSortField `xml:"field"`
}

// OverrideSortField represents a sort field in override sorting.
type OverrideSortField struct {
	Order string `xml:"order"`
}

// OverridePagination represents pagination information in the response.
type OverridePagination struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// OverrideCount represents count information in the response.
type OverrideCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
