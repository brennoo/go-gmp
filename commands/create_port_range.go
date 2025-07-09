package commands

import "encoding/xml"

// CreatePortRangeCommand represents a create_port_range command request.
type CreatePortRangeCommand struct {
	XMLName  xml.Name                `xml:"create_port_range"`
	Comment  string                  `xml:"comment,omitempty"`
	PortList CreatePortRangePortList `xml:"port_list"`
	Start    string                  `xml:"start"`
	End      string                  `xml:"end"`
	Type     string                  `xml:"type"`
}

// CreatePortRangePortList represents the port_list element with id attribute.
type CreatePortRangePortList struct {
	ID string `xml:"id,attr"`
}

// CreatePortRangeResponse represents a create_port_range command response.
type CreatePortRangeResponse struct {
	XMLName    xml.Name `xml:"create_port_range_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
