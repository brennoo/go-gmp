package commands

import "encoding/xml"

// GetPreferences represents a get_preferences command request.
type GetPreferences struct {
	XMLName    xml.Name `xml:"get_preferences"`
	NVTOID     string   `xml:"nvt_oid,attr,omitempty"`
	ConfigID   string   `xml:"config_id,attr,omitempty"`
	Preference string   `xml:"preference,attr,omitempty"`
}

// GetPreferencesResponse represents a get_preferences command response.
type GetPreferencesResponse struct {
	XMLName    xml.Name     `xml:"get_preferences_response"`
	Status     string       `xml:"status,attr"`
	StatusText string       `xml:"status_text,attr"`
	Preference []Preference `xml:"preference"`
}

// Preference represents a preference element in the get_preferences response.
type Preference struct {
	NVT     PreferenceNVT `xml:"nvt"`
	Name    string        `xml:"name"`
	ID      string        `xml:"id"`
	Type    string        `xml:"type"`
	Value   string        `xml:"value"`
	Alt     []string      `xml:"alt"`
	Default string        `xml:"default"`
}

// PreferenceNVT represents the NVT element for a preference in the get_preferences response.
type PreferenceNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name"`
}
