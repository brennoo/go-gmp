package commands

import (
	"encoding/xml"
	"time"
)

// GetTargets represents a get_targets command request.
type GetTargets struct {
	XMLName  xml.Name `xml:"get_targets"`
	TargetID string   `xml:"target_id,attr,omitempty"`
	Filter   string   `xml:"filter,attr,omitempty"`
	FiltID   string   `xml:"filt_id,attr,omitempty"`
	Trash    bool     `xml:"trash,attr,omitempty"`
	Tasks    bool     `xml:"tasks,attr,omitempty"`
	Details  bool     `xml:"details,attr,omitempty"`
}

// GetTargetsResponse represents a get_targets command response.
type GetTargetsResponse struct {
	XMLName     xml.Name                   `xml:"get_targets_response"`
	Status      string                     `xml:"status,attr"`
	StatusText  string                     `xml:"status_text,attr"`
	Target      []GetTargetsResponseTarget `xml:"target,omitempty"`
	Filters     *GetTargetsResponseFilters `xml:"filters,omitempty"`
	Sort        *GetTargetsResponseSort    `xml:"sort,omitempty"`
	Targets     *GetTargetsResponseTargets `xml:"targets,omitempty"`
	TargetCount *GetTargetsResponseCount   `xml:"target_count,omitempty"`
}

// GetTargetsResponseTarget represents a target element in the get_targets response.
type GetTargetsResponseTarget struct {
	ID                   string                         `xml:"id,attr"`
	Owner                *GetTargetsResponseTargetOwner `xml:"owner,omitempty"`
	Name                 string                         `xml:"name,omitempty"`
	Comment              string                         `xml:"comment,omitempty"`
	CreationTime         time.Time                      `xml:"creation_time,omitempty"`
	ModificationTime     time.Time                      `xml:"modification_time,omitempty"`
	Writable             bool                           `xml:"writable,omitempty"`
	InUse                bool                           `xml:"in_use,omitempty"`
	Permissions          *GetTargetsResponsePermissions `xml:"permissions,omitempty"`
	UserTags             *GetTargetsResponseUserTags    `xml:"user_tags,omitempty"`
	Hosts                string                         `xml:"hosts,omitempty"`
	ExcludeHosts         string                         `xml:"exclude_hosts,omitempty"`
	MaxHosts             int                            `xml:"max_hosts,omitempty"`
	SSHCredential        *GetTargetsResponseCredential  `xml:"ssh_credential,omitempty"`
	SMBCredential        *GetTargetsResponseCredential  `xml:"smb_credential,omitempty"`
	ESXICredential       *GetTargetsResponseCredential  `xml:"esxi_credential,omitempty"`
	Krb5Credential       *GetTargetsResponseCredential  `xml:"krb5_credential,omitempty"`
	SNMPCredential       *GetTargetsResponseCredential  `xml:"snmp_credential,omitempty"`
	SSHElevateCredential *GetTargetsResponseCredential  `xml:"ssh_elevate_credential,omitempty"`
	PortRange            *GetTargetsResponsePortRange   `xml:"port_range,omitempty"`
	PortList             *GetTargetsResponsePortList    `xml:"port_list,omitempty"`
	Tasks                *GetTargetsResponseTasks       `xml:"tasks,omitempty"`
}

type GetTargetsResponseTargetOwner struct {
	Name string `xml:"name,omitempty"`
}

type GetTargetsResponsePermissions struct {
	Permission []GetTargetsResponsePermission `xml:"permission,omitempty"`
}

type GetTargetsResponsePermission struct {
	Name string `xml:"name,omitempty"`
}

type GetTargetsResponseUserTags struct {
	Count int                             `xml:"count,attr,omitempty"`
	Tag   []GetTargetsResponseUserTagsTag `xml:"tag,omitempty"`
}

type GetTargetsResponseUserTagsTag struct {
	ID    string `xml:"id,attr,omitempty"`
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

type GetTargetsResponseCredential struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

type GetTargetsResponsePortRange struct {
	Start int `xml:"start,omitempty"`
	End   int `xml:"end,omitempty"`
}

type GetTargetsResponsePortList struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

type GetTargetsResponseTasks struct {
	Count int                      `xml:"count,attr,omitempty"`
	Task  []GetTargetsResponseTask `xml:"task,omitempty"`
}

type GetTargetsResponseTask struct {
	ID string `xml:"id,attr,omitempty"`
}

type GetTargetsResponseFilters struct {
	ID      string `xml:"id,attr,omitempty"`
	Term    string `xml:"term,omitempty"`
	Name    string `xml:"name,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

type GetTargetsResponseSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

type GetTargetsResponseTargets struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

type GetTargetsResponseCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
