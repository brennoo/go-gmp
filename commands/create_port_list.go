package commands

import "encoding/xml"

// CreatePortList is a create_port_list command request.
type CreatePortList struct {
	XMLName   xml.Name `xml:"create_port_list"`
	Name      string   `xml:"name"`
	Comment   string   `xml:"comment,omitempty"`
	Copy      string   `xml:"copy,omitempty"`
	PortRange string   `xml:"port_range"`
}

// CreatePortListResponse is a create_port_list command response.
type CreatePortListResponse struct {
	XMLName    xml.Name `xml:"create_port_list_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
