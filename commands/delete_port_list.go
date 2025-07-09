package commands

import "encoding/xml"

// DeletePortList represents a delete_port_list command request.
type DeletePortList struct {
	XMLName    xml.Name `xml:"delete_port_list"`
	PortListID string   `xml:"port_list_id,attr"`
	Ultimate   string   `xml:"ultimate,attr"` // "1" to remove entirely, "0" to move to trashcan
}

// DeletePortListResponse represents a delete_port_list command response.
type DeletePortListResponse struct {
	XMLName    xml.Name `xml:"delete_port_list_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
