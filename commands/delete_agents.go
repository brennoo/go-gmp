package commands

import "encoding/xml"

// DeleteAgents represents a delete_agents command request.
type DeleteAgents struct {
	XMLName xml.Name `xml:"delete_agents"`
	Agents  []Agent  `xml:"agents>agent"`
}

// DeleteAgentsResponse represents a delete_agents command response.
type DeleteAgentsResponse struct {
	XMLName    xml.Name `xml:"delete_agents_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteAgentsResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
