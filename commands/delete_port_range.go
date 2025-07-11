package commands

import "encoding/xml"

// DeletePortRange represents a delete_port_range command request.
type DeletePortRange struct {
	XMLName     xml.Name `xml:"delete_port_range"`
	PortRangeID string   `xml:"port_range_id,attr"`
}

// DeletePortRangeResponse represents a delete_port_range command response.
type DeletePortRangeResponse struct {
	XMLName    xml.Name `xml:"delete_port_range_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
