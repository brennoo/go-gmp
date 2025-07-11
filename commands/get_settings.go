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
	XMLName      xml.Name  `xml:"get_settings_response"`
	Status       string    `xml:"status,attr"`
	StatusText   string    `xml:"status_text,attr"`
	Settings     *Settings `xml:"settings,omitempty"`
	SettingsList []Setting `xml:"setting,omitempty"`
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
