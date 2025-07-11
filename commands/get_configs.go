package commands

import (
	"encoding/xml"
)

// GetConfigs represents a get_configs command request.
type GetConfigs struct {
	XMLName     xml.Name `xml:"get_configs"`
	ConfigID    string   `xml:"config_id,attr,omitempty"`
	Filter      string   `xml:"filter,attr,omitempty"`
	FiltID      string   `xml:"filt_id,attr,omitempty"`
	Trash       string   `xml:"trash,attr,omitempty"`
	Details     string   `xml:"details,attr,omitempty"`
	Families    string   `xml:"families,attr,omitempty"`
	Preferences string   `xml:"preferences,attr,omitempty"`
	Tasks       string   `xml:"tasks,attr,omitempty"`
	UsageType   string   `xml:"usage_type,attr,omitempty"`
}

// GetConfigsResponse represents a get_configs command response.
type GetConfigsResponse struct {
	XMLName     xml.Name      `xml:"get_configs_response"`
	Status      string        `xml:"status,attr"`
	StatusText  string        `xml:"status_text,attr"`
	Config      []Config      `xml:"config"`
	Filters     ConfigFilters `xml:"filters"`
	Sort        ConfigSort    `xml:"sort"`
	Configs     ConfigConfigs `xml:"configs"`
	ConfigCount ConfigCount   `xml:"config_count"`
}

// Config represents a config element in the get_configs response.
type Config struct {
	ID               string              `xml:"id,attr"`
	Owner            ConfigOwner         `xml:"owner"`
	Name             string              `xml:"name"`
	Comment          string              `xml:"comment"`
	CreationTime     string              `xml:"creation_time"`
	ModificationTime string              `xml:"modification_time"`
	FamilyCount      ConfigFamilyCount   `xml:"family_count"`
	NVTCount         ConfigNVTCount      `xml:"nvt_count"`
	Type             string              `xml:"type"`
	UsageType        string              `xml:"usage_type"`
	MaxNVTCount      int                 `xml:"max_nvt_count"`
	KnownNVTCount    int                 `xml:"known_nvt_count"`
	Scanner          *ConfigScanner      `xml:"scanner,omitempty"`
	InUse            string              `xml:"in_use"`
	Writable         string              `xml:"writable"`
	Permissions      ConfigPermissions   `xml:"permissions"`
	UserTags         ConfigUserTags      `xml:"user_tags"`
	Tasks            *ConfigTasks        `xml:"tasks,omitempty"`
	Families         *ConfigFamilies     `xml:"families,omitempty"`
	Preferences      *ConfigPreferences  `xml:"preferences,omitempty"`
	NVTSelectors     *ConfigNVTSelectors `xml:"nvt_selectors,omitempty"`
	Predefined       string              `xml:"predefined"`
	Deprecated       *string             `xml:"deprecated,omitempty"`
}

// ConfigScanner represents the scanner element in a config.
type ConfigScanner struct {
	ID    string `xml:"id,attr"`
	Trash string `xml:"trash"`
}

// ConfigTasks represents the tasks element in a config.
type ConfigTasks struct {
	Task []ConfigTask `xml:"task"`
}

// ConfigTask represents a task element in config tasks.
type ConfigTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
}

// ConfigFamilies represents the families element in a config.
type ConfigFamilies struct {
	Family []ConfigFamily `xml:"family"`
}

// ConfigFamily represents a family element in config families.
type ConfigFamily struct {
	Name        string `xml:"name"`
	Type        int    `xml:"type"`
	NVTCount    int    `xml:"nvt_count"`
	MaxNVTCount int    `xml:"max_nvt_count"`
	Growing     string `xml:"growing"`
}

// ConfigPreferences represents the preferences element in a config.
type ConfigPreferences struct {
	Preference []ConfigPreference `xml:"preference"`
}

// ConfigPreference represents a preference element in config preferences.
type ConfigPreference struct {
	NVT     ConfigPreferenceNVT `xml:"nvt"`
	Name    string              `xml:"name"`
	HRName  string              `xml:"hr_name,omitempty"`
	ID      string              `xml:"id"`
	Type    string              `xml:"type"`
	Value   string              `xml:"value"`
	Default string              `xml:"default,omitempty"`
	Alt     []string            `xml:"alt,omitempty"`
}

// ConfigPreferenceNVT represents the NVT element in a config preference.
type ConfigPreferenceNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
}

// ConfigNVTSelectors represents the NVT selectors element in a config.
type ConfigNVTSelectors struct {
	NVTSelector []ConfigNVTSelector `xml:"nvt_selector"`
}

// ConfigNVTSelector represents an NVT selector element in config NVT selectors.
type ConfigNVTSelector struct {
	Name        string `xml:"name"`
	Include     string `xml:"include"`
	Type        int    `xml:"type"`
	FamilyOrNVT string `xml:"family_or_nvt"`
}

// ConfigOwner represents the owner element in a config.
type ConfigOwner struct {
	Name string `xml:"name"`
}

// ConfigFamilyCount represents the family count element in a config.
type ConfigFamilyCount struct {
	Value   string `xml:",chardata"`
	Growing string `xml:"growing"`
}

// ConfigNVTCount represents the NVT count element in a config.
type ConfigNVTCount struct {
	Value   string `xml:",chardata"`
	Growing string `xml:"growing"`
}

// ConfigPermissions represents the permissions element in a config.
type ConfigPermissions struct {
	Permission []ConfigPermission `xml:"permission"`
}

// ConfigPermission represents a permission element in config permissions.
type ConfigPermission struct {
	Name string `xml:"name"`
}

// ConfigUserTags represents the user tags element in a config.
type ConfigUserTags struct {
	Count int `xml:"count"`
}

// ConfigFilters represents the filters element in the response.
type ConfigFilters struct {
	ID       string         `xml:"id,attr"`
	Term     string         `xml:"term"`
	Name     string         `xml:"name"`
	Keywords ConfigKeywords `xml:"keywords"`
}

// ConfigKeywords represents the keywords element in filters.
type ConfigKeywords struct {
	Keyword []ConfigKeyword `xml:"keyword"`
}

// ConfigKeyword represents a keyword element in filters keywords.
type ConfigKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// ConfigSort represents the sort element in the response.
type ConfigSort struct {
	Value string          `xml:",chardata"`
	Field ConfigSortField `xml:"field"`
}

// ConfigSortField represents the field element in sort.
type ConfigSortField struct {
	Order string `xml:"order"`
}

// Configs represents the configs element in the response.
type ConfigConfigs struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// ConfigCount represents the config count element in the response.
type ConfigCount struct {
	Value string `xml:",chardata"`
	Page  int    `xml:"page"`
}
