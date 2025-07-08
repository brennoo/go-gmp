package gmp

import "encoding/xml"

// CreatePortListCommand represents a create_port_list command request.
type CreatePortListCommand struct {
	XMLName   xml.Name `xml:"create_port_list"`
	Name      string   `xml:"name"`
	Comment   string   `xml:"comment,omitempty"`
	Copy      string   `xml:"copy,omitempty"`
	PortRange string   `xml:"port_range"`
}

// CreatePortListResponse represents a create_port_list command response.
type CreatePortListResponse struct {
	XMLName    xml.Name `xml:"create_port_list_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
