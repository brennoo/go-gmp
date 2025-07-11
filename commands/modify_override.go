package commands

import "encoding/xml"

// ModifyOverride represents a modify_override command request.
type ModifyOverride struct {
	XMLName     xml.Name     `xml:"modify_override"`
	OverrideID  string       `xml:"override_id,attr"`
	Active      *int         `xml:"active,omitempty"`
	NVT         *OverrideNVT `xml:"nvt,omitempty"`
	Hosts       string       `xml:"hosts,omitempty"`
	NewSeverity string       `xml:"new_severity,omitempty"`
	NewThreat   string       `xml:"new_threat,omitempty"`
	Port        string       `xml:"port,omitempty"`
	Result      string       `xml:"result,omitempty"`
	Task        string       `xml:"task,omitempty"`
	Text        string       `xml:"text"`
	Severity    string       `xml:"severity,omitempty"`
	Threat      string       `xml:"threat,omitempty"`
}

// ModifyOverrideResponse represents a modify_override command response.
type ModifyOverrideResponse struct {
	XMLName    xml.Name `xml:"modify_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
