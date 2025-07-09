package commands

import "encoding/xml"

// DeleteOverrideCommand represents a delete_override command request.
type DeleteOverrideCommand struct {
	XMLName    xml.Name `xml:"delete_override"`
	OverrideID string   `xml:"override_id,attr"`
	Ultimate   string   `xml:"ultimate,attr"`
}

// DeleteOverrideResponse represents a delete_override command response.
type DeleteOverrideResponse struct {
	XMLName    xml.Name `xml:"delete_override_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
