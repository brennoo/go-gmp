package commands

import "encoding/xml"

// CreateOverrideCommand represents a create_override command request.
type CreateOverrideCommand struct {
	XMLName     xml.Name          `xml:"create_override"`
	Text        string            `xml:"text"`
	NVT         CreateOverrideNVT `xml:"nvt"`
	Active      *int              `xml:"active,omitempty"`
	Copy        string            `xml:"copy,omitempty"`
	Hosts       string            `xml:"hosts,omitempty"`
	NewSeverity string            `xml:"new_severity,omitempty"`
	NewThreat   string            `xml:"new_threat,omitempty"`
	Port        string            `xml:"port,omitempty"`
	Result      string            `xml:"result,omitempty"`
	Severity    string            `xml:"severity,omitempty"`
	Task        string            `xml:"task,omitempty"`
	Threat      string            `xml:"threat,omitempty"`
}

// CreateOverrideNVT represents the <nvt> element for create_override.
type CreateOverrideNVT struct {
	OID string `xml:"oid,attr"`
}

// CreateOverrideResponse represents a create_override command response.
type CreateOverrideResponse struct {
	XMLName    xml.Name `xml:"create_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
