package commands

import "encoding/xml"

// ModifyPortListCommand represents a modify_port_list command request.
type ModifyPortListCommand struct {
	XMLName    xml.Name `xml:"modify_port_list"`
	PortListID string   `xml:"port_list_id,attr"`
	Name       string   `xml:"name,omitempty"`
	Comment    string   `xml:"comment,omitempty"`
}

// ModifyPortListResponse represents a modify_port_list command response.
type ModifyPortListResponse struct {
	XMLName    xml.Name `xml:"modify_port_list_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
