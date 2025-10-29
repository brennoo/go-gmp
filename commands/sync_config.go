package commands

import "encoding/xml"

// SyncConfig represents a sync_config command request.
type SyncConfig struct {
	XMLName xml.Name `xml:"sync_config"`
}

// SyncConfigResponse represents a sync_config command response.
type SyncConfigResponse struct {
	XMLName    xml.Name `xml:"sync_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *SyncConfigResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
