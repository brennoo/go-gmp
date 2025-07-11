package commands

import "encoding/xml"

// DeleteConfig represents a delete_config command request.
type DeleteConfig struct {
	XMLName  xml.Name `xml:"delete_config"`
	ConfigID string   `xml:"config_id,attr"`
	Ultimate string   `xml:"ultimate,attr"`
}

// DeleteConfigResponse represents a delete_config command response.
type DeleteConfigResponse struct {
	XMLName    xml.Name `xml:"delete_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
