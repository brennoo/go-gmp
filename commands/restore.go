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

// GetStatus returns the status and status text from the response.
func (r *RestoreResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
