package gmp

import "encoding/xml"

// DeleteConfigCommand represents a delete_config command request.
type DeleteConfigCommand struct {
	XMLName  xml.Name `xml:"delete_config"`
	ConfigID string   `xml:"config_id,attr"`
	Ultimate string   `xml:"ultimate,attr"` // "1" to remove entirely, "0" to move to trashcan
}

// DeleteConfigResponse represents a delete_config command response.
type DeleteConfigResponse struct {
	XMLName    xml.Name `xml:"delete_config_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
