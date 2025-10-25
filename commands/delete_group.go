package commands

import "encoding/xml"

// DeleteGroup represents a delete_group command request.
type DeleteGroup struct {
	XMLName  xml.Name `xml:"delete_group"`
	GroupID  string   `xml:"group_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteGroupResponse represents a delete_group command response.
type DeleteGroupResponse struct {
	XMLName    xml.Name `xml:"delete_group_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
