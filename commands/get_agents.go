package commands

import "encoding/xml"

// GetAgents represents a get_agents command request.
type GetAgents struct {
	XMLName xml.Name `xml:"get_agents"`
}

// GetAgentsResponse represents a get_agents command response.
type GetAgentsResponse struct {
	XMLName    xml.Name `xml:"get_agents_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Agents     []Agent  `xml:"agent"`
}
