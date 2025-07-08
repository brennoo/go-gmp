package gmp

import (
	"encoding/xml"
	"time"
)

// GetConfigsCommand represents a get_configs command request.
type GetConfigsCommand struct {
	XMLName     xml.Name `xml:"get_configs"`
	ConfigID    string   `xml:"config_id,attr,omitempty"`
	Filter      string   `xml:"filter,attr,omitempty"`
	FiltID      string   `xml:"filt_id,attr,omitempty"`
	Trash       bool     `xml:"trash,attr,omitempty"`
	Details     bool     `xml:"details,attr,omitempty"`
	Families    bool     `xml:"families,attr,omitempty"`
	Preferences bool     `xml:"preferences,attr,omitempty"`
	Tasks       bool     `xml:"tasks,attr,omitempty"`
	UsageType   string   `xml:"usage_type,attr,omitempty"`
}

// GetConfigsResponse represents a get_configs command response.
type GetConfigsResponse struct {
	XMLName     xml.Name                      `xml:"get_configs_response"`
	Status      string                        `xml:"status,attr"`
	StatusText  string                        `xml:"status_text,attr"`
	Config      []getConfigsResponseConfig    `xml:"config"`
	Filters     getConfigsResponseFilters     `xml:"filters"`
	Sort        getConfigsResponseSort        `xml:"sort"`
	Configs     getConfigsResponseConfigs     `xml:"configs"`
	ConfigCount getConfigsResponseConfigCount `xml:"config_count"`
}

// getConfigsResponseConfig represents a config element in the get_configs response.
type getConfigsResponseConfig struct {
	ID               string                                `xml:"id,attr"`
	Owner            getConfigsResponseConfigOwner         `xml:"owner"`
	Name             string                                `xml:"name"`
	Comment          string                                `xml:"comment"`
	CreationTime     time.Time                             `xml:"creation_time"`
	ModificationTime time.Time                             `xml:"modification_time"`
	FamilyCount      getConfigsResponseConfigFamilyCount   `xml:"family_count"`
	NVTCount         getConfigsResponseConfigNVTCount      `xml:"nvt_count"`
	Type             string                                `xml:"type"`
	UsageType        string                                `xml:"usage_type"`
	MaxNVTCount      int                                   `xml:"max_nvt_count"`
	KnownNVTCount    int                                   `xml:"known_nvt_count"`
	Scanner          *getConfigsResponseConfigScanner      `xml:"scanner,omitempty"`
	InUse            bool                                  `xml:"in_use"`
	Writable         bool                                  `xml:"writable"`
	Permissions      getConfigsResponseConfigPermissions   `xml:"permissions"`
	UserTags         getConfigsResponseConfigUserTags      `xml:"user_tags"`
	Tasks            *getConfigsResponseConfigTasks        `xml:"tasks,omitempty"`
	Families         *getConfigsResponseConfigFamilies     `xml:"families,omitempty"`
	Preferences      *getConfigsResponseConfigPreferences  `xml:"preferences,omitempty"`
	NVTSelectors     *getConfigsResponseConfigNVTSelectors `xml:"nvt_selectors,omitempty"`
	Predefined       bool                                  `xml:"predefined"`
	Deprecated       *bool                                 `xml:"deprecated,omitempty"`
}

type getConfigsResponseConfigScanner struct {
	ID    string `xml:"id,attr"`
	Trash bool   `xml:"trash"`
}

type getConfigsResponseConfigTasks struct {
	Task []getConfigsResponseConfigTask `xml:"task"`
}

type getConfigsResponseConfigTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
}

type getConfigsResponseConfigFamilies struct {
	Family []getConfigsResponseConfigFamily `xml:"family"`
}

type getConfigsResponseConfigFamily struct {
	Name        string `xml:"name"`
	Type        int    `xml:"type"`
	NVTCount    int    `xml:"nvt_count"`
	MaxNVTCount int    `xml:"max_nvt_count"`
	Growing     bool   `xml:"growing"`
}

type getConfigsResponseConfigPreferences struct {
	Preference []getConfigsResponseConfigPreference `xml:"preference"`
}

type getConfigsResponseConfigPreference struct {
	NVT     getConfigsResponseConfigPreferenceNVT `xml:"nvt"`
	Name    string                                `xml:"name"`
	HRName  string                                `xml:"hr_name,omitempty"`
	ID      string                                `xml:"id"`
	Type    string                                `xml:"type"`
	Value   string                                `xml:"value"`
	Default string                                `xml:"default,omitempty"`
	Alt     []string                              `xml:"alt,omitempty"`
}

type getConfigsResponseConfigPreferenceNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
}

type getConfigsResponseConfigNVTSelectors struct {
	NVTSelector []getConfigsResponseConfigNVTSelector `xml:"nvt_selector"`
}

type getConfigsResponseConfigNVTSelector struct {
	Name        string `xml:"name"`
	Include     bool   `xml:"include"`
	Type        int    `xml:"type"`
	FamilyOrNVT string `xml:"family_or_nvt"`
}

type getConfigsResponseConfigOwner struct {
	Name string `xml:"name"`
}

type getConfigsResponseConfigFamilyCount struct {
	Value   string `xml:",chardata"`
	Growing string `xml:"growing"`
}

type getConfigsResponseConfigNVTCount struct {
	Value   string `xml:",chardata"`
	Growing string `xml:"growing"`
}

type getConfigsResponseConfigPermissions struct {
	Permission []getConfigsResponseConfigPermissionsPermission `xml:"permission"`
}

type getConfigsResponseConfigPermissionsPermission struct {
	Name string `xml:"name"`
}

type getConfigsResponseConfigUserTags struct {
	Count int `xml:"count"`
}

type getConfigsResponseFilters struct {
	ID       string                            `xml:"id,attr"`
	Term     string                            `xml:"term"`
	Name     string                            `xml:"name"`
	Keywords getConfigsResponseFiltersKeywords `xml:"keywords"`
}

type getConfigsResponseFiltersKeywords struct {
	Keyword []getConfigsResponseFiltersKeywordsKeyword `xml:"keyword"`
}

type getConfigsResponseFiltersKeywordsKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

type getConfigsResponseSort struct {
	Value string                      `xml:",chardata"`
	Field getConfigsResponseSortField `xml:"field"`
}

type getConfigsResponseSortField struct {
	Value string `xml:",chardata"`
	Order string `xml:"order"`
}

type getConfigsResponseConfigs struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

type getConfigsResponseConfigCount struct {
	Value    string `xml:",chardata"`
	Filtered int    `xml:"filtered"`
	Page     int    `xml:"page"`
}
