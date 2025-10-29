package commands

import "encoding/xml"

// GetSettings represents a get_settings command request.
type GetSettings struct {
	XMLName   xml.Name `xml:"get_settings"`
	SettingID string   `xml:"setting_id,attr,omitempty"`
	Filter    string   `xml:"filter,attr,omitempty"`
	First     int      `xml:"first,attr,omitempty"`
	Max       int      `xml:"max,attr,omitempty"`
	SortOrder string   `xml:"sort_order,attr,omitempty"`
	SortField string   `xml:"sort_field,attr,omitempty"`
}

// GetSettingsResponse represents a get_settings command response.
type GetSettingsResponse struct {
	XMLName      xml.Name        `xml:"get_settings_response"`
	Status       string          `xml:"status,attr"`
	StatusText   string          `xml:"status_text,attr"`
	Filters      *SettingFilters `xml:"filters,omitempty"`
	Settings     *Settings       `xml:"settings,omitempty"`
	SettingsList []Setting       `xml:"setting,omitempty"`
	SettingCount *SettingCount   `xml:"setting_count,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetSettingsResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Settings represents a <settings> element in the get_settings response.
type Settings struct {
	Start   int       `xml:"start,attr,omitempty"`
	Max     int       `xml:"max,attr,omitempty"`
	Setting []Setting `xml:"setting,omitempty"`
}

// Setting represents a <setting> element in the get_settings response.
type Setting struct {
	ID      string `xml:"id,attr,omitempty"`
	Name    string `xml:"name,omitempty"`
	Comment string `xml:"comment,omitempty"`
	Value   string `xml:"value,omitempty"`
}

// SettingFilters represents the filters element in the response.
type SettingFilters struct {
	Term string `xml:"term,omitempty"`
}

// SettingCount represents the setting_count element in the response.
type SettingCount struct {
	Filtered int `xml:"filtered,omitempty"`
	Page     int `xml:"page,omitempty"`
}
