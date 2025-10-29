package commands

import "encoding/xml"

// DeleteConfig represents a delete_config command request.
type DeleteConfig struct {
	XMLName  xml.Name `xml:"delete_config"`
	ConfigID string   `xml:"config_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteConfigResponse represents a delete_config command response.
type DeleteConfigResponse struct {
	XMLName    xml.Name `xml:"delete_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteConfigResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
