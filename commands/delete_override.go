package commands

import "encoding/xml"

// DeleteOverride represents a delete_override command request.
type DeleteOverride struct {
	XMLName    xml.Name `xml:"delete_override"`
	OverrideID string   `xml:"override_id,attr"`
	Ultimate   bool     `xml:"ultimate,attr"`
}

// DeleteOverrideResponse represents a delete_override command response.
type DeleteOverrideResponse struct {
	XMLName    xml.Name `xml:"delete_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteOverrideResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
