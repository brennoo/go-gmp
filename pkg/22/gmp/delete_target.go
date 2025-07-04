package gmp

import (
	"encoding/xml"
)

type DeleteTargetCommand struct {
	XMLName  xml.Name `xml:"delete_target"`
	TargetID string   `xml:"target_id,attr,omitempty"`
	Ultimate bool     `xml:"ultimate,attr,omitempty"`
}

type DeleteTargetResponse struct {
	XMLName    xml.Name `xml:"delete_target_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
