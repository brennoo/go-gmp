package commands

import "encoding/xml"

// ModifyAgents represents a modify_agents command request.
type ModifyAgents struct {
	XMLName           xml.Name `xml:"modify_agents"`
	Agents            []Agent  `xml:"agents>agent"`
	Authorized        int      `xml:"authorized,omitempty"`
	MinInterval       int      `xml:"min_interval,omitempty"`
	HeartbeatInterval int      `xml:"heartbeat_interval,omitempty"`
	Schedule          string   `xml:"schedule,omitempty"`
	Comment           string   `xml:"comment,omitempty"`
}

// ModifyAgentsResponse represents a modify_agents command response.
type ModifyAgentsResponse struct {
	XMLName    xml.Name `xml:"modify_agents_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
