package commands

import "encoding/xml"

// SyncConfigCommand represents a sync_config command request.
type SyncConfigCommand struct {
	XMLName xml.Name `xml:"sync_config"`
}

// SyncConfigResponse represents a sync_config command response.
type SyncConfigResponse struct {
	XMLName    xml.Name `xml:"sync_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
