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

// Target represents a <target> element in the get_targets response.
type Target struct {
	ID                   string             `xml:"id,attr"`
	Owner                *TargetOwner       `xml:"owner,omitempty"`
	Name                 string             `xml:"name,omitempty"`
	Comment              string             `xml:"comment,omitempty"`
	CreationTime         time.Time          `xml:"creation_time,omitempty"`
	ModificationTime     time.Time          `xml:"modification_time,omitempty"`
	Writable             bool               `xml:"writable,omitempty"`
	InUse                bool               `xml:"in_use,omitempty"`
	Permissions          *TargetPermissions `xml:"permissions,omitempty"`
	UserTags             *TargetUserTags    `xml:"user_tags,omitempty"`
	Hosts                string             `xml:"hosts,omitempty"`
	ExcludeHosts         string             `xml:"exclude_hosts,omitempty"`
	MaxHosts             int                `xml:"max_hosts,omitempty"`
	SSHCredential        *TargetCredential  `xml:"ssh_credential,omitempty"`
	SMBCredential        *TargetCredential  `xml:"smb_credential,omitempty"`
	ESXICredential       *TargetCredential  `xml:"esxi_credential,omitempty"`
	Krb5Credential       *TargetCredential  `xml:"krb5_credential,omitempty"`
	SNMPCredential       *TargetCredential  `xml:"snmp_credential,omitempty"`
	SSHElevateCredential *TargetCredential  `xml:"ssh_elevate_credential,omitempty"`
	PortRange            *TargetPortRange   `xml:"port_range,omitempty"`
	PortList             *TargetPortList    `xml:"port_list,omitempty"`
	AliveTests           string             `xml:"alive_tests,omitempty"`
	ReverseLookupOnly    bool               `xml:"reverse_lookup_only,omitempty"`
	ReverseLookupUnify   bool               `xml:"reverse_lookup_unify,omitempty"`
	AllowSimultaneousIPs bool               `xml:"allow_simultaneous_ips,omitempty"`
	Tasks                *TargetTasks       `xml:"tasks,omitempty"`
}

// TargetOwner represents the <owner> element for a target.
type TargetOwner struct {
	Name string `xml:"name,omitempty"`
}

// TargetPermissions represents the <permissions> element for a target.
type TargetPermissions struct {
	Permission []TargetPermission `xml:"permission,omitempty"`
}

// TargetPermission represents a <permission> element for a target.
type TargetPermission struct {
	Name string `xml:"name,omitempty"`
}

// TargetUserTags represents the <user_tags> element for a target.
type TargetUserTags struct {
	Count int             `xml:"count,attr,omitempty"`
	Tag   []TargetUserTag `xml:"tag,omitempty"`
}

// TargetUserTag represents a <tag> element in user_tags.
type TargetUserTag struct {
	ID    string `xml:"id,attr,omitempty"`
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

// TargetCredential represents a credential element for a target or create_target. The Port field is only used for credentials that support it (e.g., ssh_credential, ssh_lsc_credential).
type TargetCredential struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
	Port string `xml:"port,omitempty"`
}

// TargetPortRange represents the <port_range> element for a target.
type TargetPortRange struct {
	Start int `xml:"start,omitempty"`
	End   int `xml:"end,omitempty"`
}

// TargetPortList represents the <port_list> element for a target.
type TargetPortList struct {
	ID   string `xml:"id,attr,omitempty"`
	Name string `xml:"name,omitempty"`
}

// TargetTasks represents the <tasks> element for a target.
type TargetTasks struct {
	Count int          `xml:"count,attr,omitempty"`
	Task  []TargetTask `xml:"task,omitempty"`
}

// TargetTask represents a <task> element in tasks.
type TargetTask struct {
	ID string `xml:"id,attr,omitempty"`
}

// TargetFilters represents the <filters> element in the response.
type TargetFilters struct {
	ID      string `xml:"id,attr,omitempty"`
	Term    string `xml:"term,omitempty"`
	Name    string `xml:"name,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// TargetSort represents the <sort> element in the response.
type TargetSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// TargetTargets represents the <targets> meta element in the response.
type TargetTargets struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// TargetCount represents the <target_count> element in the response.
type TargetCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}

// GetTargetsResponse represents a get_targets command response.
type GetTargetsResponse struct {
	XMLName     xml.Name       `xml:"get_targets_response"`
	Status      string         `xml:"status,attr"`
	StatusText  string         `xml:"status_text,attr"`
	Target      []Target       `xml:"target,omitempty"`
	Filters     *TargetFilters `xml:"filters,omitempty"`
	Sort        *TargetSort    `xml:"sort,omitempty"`
	Targets     *TargetTargets `xml:"targets,omitempty"`
	TargetCount *TargetCount   `xml:"target_count,omitempty"`
	Truncated   string         `xml:"truncated,omitempty"`
}
