package gmp

import "encoding/xml"

// DeleteAgentsCommand represents a delete_agents command request.
type DeleteAgentsCommand struct {
	XMLName xml.Name `xml:"delete_agents"`
	Agents  []Agent  `xml:"agents>agent"`
}

// DeleteAgentsResponse represents a delete_agents command response.
type DeleteAgentsResponse struct {
	XMLName    xml.Name `xml:"delete_agents_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
