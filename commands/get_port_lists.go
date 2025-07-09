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
	Details    bool     `xml:"details,attr,omitempty"`
	Targets    bool     `xml:"targets,attr,omitempty"`
	Trash      bool     `xml:"trash,attr,omitempty"`
}

// GetPortListsResponse represents a get_port_lists command response.
type GetPortListsResponse struct {
	XMLName       xml.Name                       `xml:"get_port_lists_response"`
	Status        string                         `xml:"status,attr"`
	StatusText    string                         `xml:"status_text,attr"`
	PortList      []GetPortListsResponsePortList `xml:"port_list,omitempty"`
	Filters       *GetPortListsResponseFilters   `xml:"filters,omitempty"`
	Sort          *GetPortListsResponseSort      `xml:"sort,omitempty"`
	PortLists     *GetPortListsResponsePortLists `xml:"port_lists,omitempty"`
	PortListCount *GetPortListsResponseCount     `xml:"port_list_count,omitempty"`
}

// GetPortListsResponsePortList represents a port_list element in the get_port_lists response.
type GetPortListsResponsePortList struct {
	ID               string                             `xml:"id,attr"`
	Owner            *GetPortListsResponsePortListOwner `xml:"owner,omitempty"`
	Name             string                             `xml:"name,omitempty"`
	Comment          string                             `xml:"comment,omitempty"`
	CreationTime     time.Time                          `xml:"creation_time,omitempty"`
	ModificationTime time.Time                          `xml:"modification_time,omitempty"`
	Writable         bool                               `xml:"writable,omitempty"`
	InUse            bool                               `xml:"in_use,omitempty"`
	Permissions      *GetPortListsResponsePermissions   `xml:"permissions,omitempty"`
	UserTags         *GetPortListsResponseUserTags      `xml:"user_tags,omitempty"`
	PortCount        *GetPortListsResponsePortCount     `xml:"port_count,omitempty"`
	PortRanges       *GetPortListsResponsePortRanges    `xml:"port_ranges,omitempty"`
	Targets          *GetPortListsResponseTargets       `xml:"targets,omitempty"`
	Predefined       bool                               `xml:"predefined,omitempty"`
	Deprecated       bool                               `xml:"deprecated,omitempty"`
}

// GetPortListsResponsePortListOwner represents the owner of a port list.
type GetPortListsResponsePortListOwner struct {
	Name string `xml:"name,omitempty"`
}

// GetPortListsResponsePermissions represents permissions for a port list.
type GetPortListsResponsePermissions struct {
	Permission []GetPortListsResponsePermission `xml:"permission,omitempty"`
}

// GetPortListsResponsePermission represents a permission for a port list.
type GetPortListsResponsePermission struct {
	Name string `xml:"name,omitempty"`
}

// GetPortListsResponseUserTags represents user tags for a port list.
type GetPortListsResponseUserTags struct {
	Count int                               `xml:"count,attr,omitempty"`
	Tag   []GetPortListsResponseUserTagsTag `xml:"tag,omitempty"`
}

// GetPortListsResponseUserTagsTag represents a user tag for a port list.
type GetPortListsResponseUserTagsTag struct {
	ID      string `xml:"id,attr,omitempty"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// GetPortListsResponsePortCount represents the port count for a port list.
type GetPortListsResponsePortCount struct {
	All int `xml:"all,omitempty"`
	TCP int `xml:"tcp,omitempty"`
	UDP int `xml:"udp,omitempty"`
}

// GetPortListsResponsePortRanges represents port ranges for a port list.
type GetPortListsResponsePortRanges struct {
	PortRange []GetPortListsResponsePortRange `xml:"port_range,omitempty"`
}

// GetPortListsResponsePortRange represents a port range for a port list.
type GetPortListsResponsePortRange struct {
	ID      string `xml:"id,attr,omitempty"`
	Start   int    `xml:"start,omitempty"`
	End     int    `xml:"end,omitempty"`
	Type    string `xml:"type,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// GetPortListsResponseTargets represents targets for a port list.
type GetPortListsResponseTargets struct {
	Target []GetPortListsResponseTarget `xml:"target,omitempty"`
}

// GetPortListsResponseTarget represents a target for a port list.
type GetPortListsResponseTarget struct {
	ID          string `xml:"id,attr,omitempty"`
	Name        string `xml:"name,omitempty"`
	Permissions string `xml:"permissions,omitempty"`
}

// GetPortListsResponseFilters represents filters for the get_port_lists response.
type GetPortListsResponseFilters struct {
	ID       string `xml:"id,attr,omitempty"`
	Term     string `xml:"term,omitempty"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

// GetPortListsResponseSort represents sorting information for the get_port_lists response.
type GetPortListsResponseSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// GetPortListsResponsePortLists represents the port_lists element in the get_port_lists response.
type GetPortListsResponsePortLists struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// GetPortListsResponseCount represents the port_list_count element in the get_port_lists response.
type GetPortListsResponseCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
