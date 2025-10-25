package commands

import (
	"encoding/xml"
	"time"
)

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

// Agent represents an agent element in the get_agents response.
type Agent struct {
	ID                string            `xml:"id,attr"`
	Owner             *AgentOwner       `xml:"owner,omitempty"`
	Name              string            `xml:"name"`
	Comment           string            `xml:"comment,omitempty"`
	CreationTime      time.Time         `xml:"creation_time,omitempty"`
	ModificationTime  time.Time         `xml:"modification_time,omitempty"`
	Writable          bool              `xml:"writable,omitempty"`
	InUse             bool              `xml:"in_use,omitempty"`
	Permissions       *AgentPermissions `xml:"permissions,omitempty"`
	Hostname          string            `xml:"hostname,omitempty"`
	AgentID           string            `xml:"agent_id,omitempty"`
	Authorized        int               `xml:"authorized,omitempty"`
	MinInterval       int               `xml:"min_interval,omitempty"`
	HeartbeatInterval int               `xml:"heartbeat_interval,omitempty"`
	ConnectionStatus  string            `xml:"connection_status,omitempty"`
	LastUpdate        time.Time         `xml:"last_update,omitempty"`
	Schedule          string            `xml:"schedule,omitempty"`
	Scanner           *AgentScanner     `xml:"scanner,omitempty"`
	IP                []string          `xml:"ip,omitempty"`
}

type AgentOwner struct {
	Name string `xml:"name"`
}

type AgentPermissions struct {
	Permission []AgentPermission `xml:"permission"`
}

type AgentPermission struct {
	Name string `xml:"name"`
}

type AgentScanner struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}
