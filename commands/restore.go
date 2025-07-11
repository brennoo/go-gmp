package commands

import "encoding/xml"

// Restore represents a restore command request.
type Restore struct {
	XMLName xml.Name `xml:"restore"`
	ID      string   `xml:"id,attr"`
}

// RestoreResponse represents a restore command response.
type RestoreResponse struct {
	XMLName    xml.Name `xml:"restore_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
