package commands

import "encoding/xml"

// ModifyOverrideCommand represents a modify_override command request.
type ModifyOverrideCommand struct {
	XMLName     xml.Name           `xml:"modify_override"`
	OverrideID  string             `xml:"override_id,attr"`
	Active      *int               `xml:"active,omitempty"`
	NVT         *ModifyOverrideNVT `xml:"nvt,omitempty"`
	Hosts       string             `xml:"hosts,omitempty"`
	NewSeverity string             `xml:"new_severity,omitempty"`
	NewThreat   string             `xml:"new_threat,omitempty"`
	Port        string             `xml:"port,omitempty"`
	Result      string             `xml:"result,omitempty"`
	Task        string             `xml:"task,omitempty"`
	Text        string             `xml:"text"`
	Severity    string             `xml:"severity,omitempty"`
	Threat      string             `xml:"threat,omitempty"`
}

// ModifyOverrideNVT represents the NVT element in modify_override.
type ModifyOverrideNVT struct {
	OID string `xml:"oid,attr"`
}

// ModifyOverrideResponse represents a modify_override command response.
type ModifyOverrideResponse struct {
	XMLName    xml.Name `xml:"modify_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
