package commands

import (
	"encoding/xml"
	"time"
)

// GetPortLists represents a get_port_lists command request.
type GetPortLists struct {
	XMLName    xml.Name `xml:"get_port_lists"`
	PortListID string   `xml:"port_list_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
	FiltID     string   `xml:"filt_id,attr,omitempty"`
	Details    string   `xml:"details,attr,omitempty"`
	Targets    string   `xml:"targets,attr,omitempty"`
	Trash      string   `xml:"trash,attr,omitempty"`
}

// GetPortListsResponse represents a get_port_lists command response.
type GetPortListsResponse struct {
	XMLName    xml.Name            `xml:"get_port_lists_response"`
	Status     string              `xml:"status,attr"`
	StatusText string              `xml:"status_text,attr"`
	PortLists  []PortList          `xml:"port_list,omitempty"`
	Filters    *PortListFilters    `xml:"filters,omitempty"`
	Sort       *PortListSort       `xml:"sort,omitempty"`
	Pagination *PortListPagination `xml:"port_lists,omitempty"`
	Count      *PortListCount      `xml:"port_list_count,omitempty"`
}

// PortList represents a port_list element in the response.
type PortList struct {
	ID               string               `xml:"id,attr"`
	Owner            *PortListOwner       `xml:"owner,omitempty"`
	Name             string               `xml:"name,omitempty"`
	Comment          string               `xml:"comment,omitempty"`
	CreationTime     time.Time            `xml:"creation_time,omitempty"`
	ModificationTime time.Time            `xml:"modification_time,omitempty"`
	Writable         bool                 `xml:"writable,omitempty"`
	InUse            bool                 `xml:"in_use,omitempty"`
	Permissions      *PortListPermissions `xml:"permissions,omitempty"`
	UserTags         *PortListUserTags    `xml:"user_tags,omitempty"`
	PortCount        *PortListPortCount   `xml:"port_count,omitempty"`
	PortRanges       *PortListPortRanges  `xml:"port_ranges,omitempty"`
	Targets          *PortListTargets     `xml:"targets,omitempty"`
	Predefined       bool                 `xml:"predefined,omitempty"`
	Deprecated       bool                 `xml:"deprecated,omitempty"`
}

// PortListOwner represents the owner of a port list.
type PortListOwner struct {
	Name string `xml:"name,omitempty"`
}

// PortListPermissions represents permissions for a port list.
type PortListPermissions struct {
	Permission []PortListPermission `xml:"permission,omitempty"`
}

// PortListPermission represents a single permission for a port list.
type PortListPermission struct {
	Name string `xml:"name,omitempty"`
}

// PortListUserTags represents user tags for a port list.
type PortListUserTags struct {
	Count int               `xml:"count,attr,omitempty"`
	Tag   []PortListUserTag `xml:"tag,omitempty"`
}

// PortListUserTag represents a user tag for a port list.
type PortListUserTag struct {
	ID      string `xml:"id,attr,omitempty"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// PortListPortCount represents the port count for a port list.
type PortListPortCount struct {
	All int `xml:"all,omitempty"`
	TCP int `xml:"tcp,omitempty"`
	UDP int `xml:"udp,omitempty"`
}

// PortListPortRanges represents port ranges for a port list.
type PortListPortRanges struct {
	PortRange []PortListPortRange `xml:"port_range,omitempty"`
}

// PortListPortRange represents a port range for a port list.
type PortListPortRange struct {
	ID      string `xml:"id,attr,omitempty"`
	Start   int    `xml:"start,omitempty"`
	End     int    `xml:"end,omitempty"`
	Type    string `xml:"type,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// PortListTargets represents targets for a port list.
type PortListTargets struct {
	Target []PortListTarget `xml:"target,omitempty"`
}

// PortListTarget represents a target for a port list.
type PortListTarget struct {
	ID          string `xml:"id,attr,omitempty"`
	Name        string `xml:"name,omitempty"`
	Permissions string `xml:"permissions,omitempty"`
}

// PortListFilters represents filters for the get_port_lists response.
type PortListFilters struct {
	ID       string `xml:"id,attr,omitempty"`
	Term     string `xml:"term,omitempty"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

// PortListSort represents sorting information for the get_port_lists response.
type PortListSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// PortListPagination represents pagination information for the get_port_lists response.
type PortListPagination struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// PortListCount represents the port_list_count element in the get_port_lists response.
type PortListCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
