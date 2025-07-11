package commands

import (
	"encoding/xml"
)

// DeleteTarget represents a delete_target command request.
type DeleteTarget struct {
	XMLName  xml.Name `xml:"delete_target"`
	TargetID string   `xml:"target_id,attr,omitempty"`
	Ultimate bool     `xml:"ultimate,attr,omitempty"`
}

// DeleteTargetResponse represents a delete_target command response.
type DeleteTargetResponse struct {
	XMLName    xml.Name `xml:"delete_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
