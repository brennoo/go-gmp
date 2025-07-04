package gmp

import (
	"encoding/xml"
	"time"
)

type GetPortListsCommand struct {
	XMLName    xml.Name `xml:"get_port_lists"`
	PortListID string   `xml:"port_list_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
	FiltID     string   `xml:"filt_id,attr,omitempty"`
	Details    bool     `xml:"details,attr,omitempty"`
	Targets    bool     `xml:"targets,attr,omitempty"`
	Trash      bool     `xml:"trash,attr,omitempty"`
}

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

type GetPortListsResponsePortListOwner struct {
	Name string `xml:"name,omitempty"`
}

type GetPortListsResponsePermissions struct {
	Permission []GetPortListsResponsePermission `xml:"permission,omitempty"`
}

type GetPortListsResponsePermission struct {
	Name string `xml:"name,omitempty"`
}

type GetPortListsResponseUserTags struct {
	Count int                               `xml:"count,attr,omitempty"`
	Tag   []GetPortListsResponseUserTagsTag `xml:"tag,omitempty"`
}

type GetPortListsResponseUserTagsTag struct {
	ID      string `xml:"id,attr,omitempty"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

type GetPortListsResponsePortCount struct {
	All int `xml:"all,omitempty"`
	TCP int `xml:"tcp,omitempty"`
	UDP int `xml:"udp,omitempty"`
}

type GetPortListsResponsePortRanges struct {
	PortRange []GetPortListsResponsePortRange `xml:"port_range,omitempty"`
}

type GetPortListsResponsePortRange struct {
	ID      string `xml:"id,attr,omitempty"`
	Start   int    `xml:"start,omitempty"`
	End     int    `xml:"end,omitempty"`
	Type    string `xml:"type,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

type GetPortListsResponseTargets struct {
	Target []GetPortListsResponseTarget `xml:"target,omitempty"`
}

type GetPortListsResponseTarget struct {
	ID          string `xml:"id,attr,omitempty"`
	Name        string `xml:"name,omitempty"`
	Permissions string `xml:"permissions,omitempty"`
}

type GetPortListsResponseFilters struct {
	ID       string `xml:"id,attr,omitempty"`
	Term     string `xml:"term,omitempty"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

type GetPortListsResponseSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

type GetPortListsResponsePortLists struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

type GetPortListsResponseCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
