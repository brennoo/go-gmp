package commands

import "encoding/xml"

// CreateOverride is a create_override command request.
type CreateOverride struct {
	XMLName     xml.Name     `xml:"create_override"`
	Text        string       `xml:"text"`
	NVT         *OverrideNVT `xml:"nvt"`
	Active      *int         `xml:"active,omitempty"`
	Copy        string       `xml:"copy,omitempty"`
	Hosts       string       `xml:"hosts,omitempty"`
	NewSeverity float64      `xml:"new_severity,omitempty"`
	NewThreat   string       `xml:"new_threat,omitempty"`
	Port        string       `xml:"port,omitempty"`
	Result      string       `xml:"result,omitempty"`
	Severity    float64      `xml:"severity,omitempty"`
	Task        string       `xml:"task,omitempty"`
	Threat      string       `xml:"threat,omitempty"`
}

// CreateOverrideResponse is a create_override command response.
type CreateOverrideResponse struct {
	XMLName    xml.Name `xml:"create_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateOverrideResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
