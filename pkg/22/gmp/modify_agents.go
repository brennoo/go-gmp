package gmp

import "encoding/xml"

// ModifyAgentsCommand represents a modify_agents command request.
type ModifyAgentsCommand struct {
	XMLName           xml.Name `xml:"modify_agents"`
	Agents            []Agent  `xml:"agents>agent"`
	Authorized        string   `xml:"authorized,omitempty"`
	MinInterval       string   `xml:"min_interval,omitempty"`
	HeartbeatInterval string   `xml:"heartbeat_interval,omitempty"`
	Schedule          string   `xml:"schedule,omitempty"`
	Comment           string   `xml:"comment,omitempty"`
}

type Agent struct {
	ID string `xml:"id,attr"`
}

// ModifyAgentsResponse represents a modify_agents command response.
type ModifyAgentsResponse struct {
	XMLName    xml.Name `xml:"modify_agents_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
