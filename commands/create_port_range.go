package commands

import "encoding/xml"

// CreatePortRange is a create_port_range command request.
type CreatePortRange struct {
	XMLName    xml.Name `xml:"create_port_range"`
	Comment    string   `xml:"comment,omitempty"`
	PortListID string   `xml:"port_list"`
	Start      string   `xml:"start"`
	End        string   `xml:"end"`
	Type       string   `xml:"type"`
}

// CreatePortRangeResponse is a create_port_range command response.
type CreatePortRangeResponse struct {
	XMLName    xml.Name `xml:"create_port_range_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
