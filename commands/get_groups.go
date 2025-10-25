package commands

import (
	"encoding/xml"
	"time"
)

// GetGroups represents a get_groups command request.
type GetGroups struct {
	XMLName xml.Name `xml:"get_groups"`
	GroupID string   `xml:"group_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
	FiltID  string   `xml:"filt_id,attr,omitempty"`
	Trash   bool     `xml:"trash,attr,omitempty"`
}

// GetGroupsResponse represents a get_groups command response.
type GetGroupsResponse struct {
	XMLName    xml.Name     `xml:"get_groups_response"`
	Status     string       `xml:"status,attr"`
	StatusText string       `xml:"status_text,attr"`
	Groups     []Group      `xml:"group"`
	Filters    GroupFilters `xml:"filters"`
	Sort       GroupSort    `xml:"sort"`
	GroupsInfo GroupInfo    `xml:"groups"`
	GroupCount GroupCount   `xml:"group_count"`
}

// Group represents a <group> element in the get_groups response.
type Group struct {
	ID               string           `xml:"id,attr"`
	Owner            GroupOwner       `xml:"owner"`
	Name             string           `xml:"name"`
	Comment          string           `xml:"comment"`
	CreationTime     time.Time        `xml:"creation_time"`
	ModificationTime time.Time        `xml:"modification_time"`
	Writable         bool             `xml:"writable"`
	InUse            bool             `xml:"in_use"`
	Permissions      GroupPermissions `xml:"permissions"`
	UserTags         *GroupUserTags   `xml:"user_tags,omitempty"`
	Users            string           `xml:"users"`
}

// GroupOwner represents the owner element in a group.
type GroupOwner struct {
	Name string `xml:"name"`
}

// GroupPermissions represents the permissions element in a group.
type GroupPermissions struct {
	Permission []GroupPermission `xml:"permission"`
}

// GroupPermission represents a permission element in group permissions.
type GroupPermission struct {
	Name string `xml:"name"`
}

// GroupUserTags represents the user tags element in a group.
type GroupUserTags struct {
	Count int            `xml:"count"`
	Tags  []GroupUserTag `xml:"tag,omitempty"`
}

// GroupUserTag represents a tag element in group user tags.
type GroupUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// GroupFilters represents the filters element in the response.
type GroupFilters struct {
	ID       string        `xml:"id,attr"`
	Term     string        `xml:"term"`
	Name     string        `xml:"name"`
	Keywords GroupKeywords `xml:"keywords"`
}

// GroupKeywords represents the keywords element in filters.
type GroupKeywords struct {
	Keyword []GroupKeyword `xml:"keyword"`
}

// GroupKeyword represents a keyword element in filters keywords.
type GroupKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// GroupSort represents the sort element in the response.
type GroupSort struct {
	Value string         `xml:",chardata"`
	Field GroupSortField `xml:"field"`
}

// GroupSortField represents the field element in sort.
type GroupSortField struct {
	Order string `xml:"order"`
}

// GroupInfo represents the groups element in the response.
type GroupInfo struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// GroupCount represents the group count element in the response.
type GroupCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
