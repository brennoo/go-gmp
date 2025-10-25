package commands

import (
	"encoding/xml"
	"time"
)

// GetUsers represents a get_users command request.
type GetUsers struct {
	XMLName xml.Name `xml:"get_users"`
	UserID  string   `xml:"user_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
}

// GetUsersResponse represents a get_users command response.
type GetUsersResponse struct {
	XMLName    xml.Name     `xml:"get_users_response"`
	Status     string       `xml:"status,attr"`
	StatusText string       `xml:"status_text,attr"`
	Users      []User       `xml:"user,omitempty"`
	Filters    *UserFilters `xml:"filters,omitempty"`
	Sort       *UserSort    `xml:"sort,omitempty"`
	UsersInfo  *UserUsers   `xml:"users,omitempty"`
	UserCount  *UserCount   `xml:"user_count,omitempty"`
	Truncated  string       `xml:"truncated,omitempty"`
}

// User represents a <user> element in the get_users response.
type User struct {
	ID               string           `xml:"id,attr,omitempty"`
	Owner            *UserOwner       `xml:"owner,omitempty"`
	Name             string           `xml:"name"`
	Comment          string           `xml:"comment,omitempty"`
	CreationTime     time.Time        `xml:"creation_time,omitempty"`
	ModificationTime time.Time        `xml:"modification_time,omitempty"`
	Writable         bool             `xml:"writable,omitempty"`
	InUse            bool             `xml:"in_use,omitempty"`
	Role             *UserRole        `xml:"role,omitempty"`
	Groups           *UserGroups      `xml:"groups,omitempty"`
	Hosts            *UserHosts       `xml:"hosts,omitempty"`
	Permissions      *UserPermissions `xml:"permissions,omitempty"`
	UserTags         *UserTags        `xml:"user_tags,omitempty"`
	Sources          *UserSources     `xml:"sources,omitempty"`
}

// UserRole represents the <role> element in a user.
type UserRole struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// UserGroups represents the <groups> element in a user.
type UserGroups struct {
	Group []UserGroup `xml:"group"`
}

// UserGroup represents a <group> element in groups.
type UserGroup struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// UserHosts represents the <hosts> element in a user.
type UserHosts struct {
	Allow string `xml:"allow,attr"`
	Value string `xml:",chardata"`
}

// UserPermissions represents the <permissions> element in a user.
type UserPermissions struct {
	Permission []UserPermission `xml:"permission"`
}

// UserPermission represents a <permission> element in permissions.
type UserPermission struct {
	Name string `xml:"name"`
}

// UserTags represents the <user_tags> element in a user.
type UserTags struct {
	Count int       `xml:"count"`
	Tag   []UserTag `xml:"tag"`
}

// UserTag represents a <tag> element in user_tags.
type UserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// UserOwner represents the <owner> element for a user.
type UserOwner struct {
	Name string `xml:"name"`
}

// UserSources represents the <sources> element in a user.
type UserSources struct {
	Source []string `xml:"source"`
}

// UserFilters represents the <filters> element in the response.
type UserFilters struct {
	ID       string               `xml:"id,attr,omitempty"`
	Term     string               `xml:"term,omitempty"`
	Name     string               `xml:"name,omitempty"`
	Comment  string               `xml:"comment,omitempty"`
	Keywords *UserFiltersKeywords `xml:"keywords,omitempty"`
}

// UserFiltersKeywords represents the <keywords> element in filters.
type UserFiltersKeywords struct {
	Keyword []UserFiltersKeywordsKeyword `xml:"keyword,omitempty"`
}

// UserFiltersKeywordsKeyword represents a <keyword> element in filters keywords.
type UserFiltersKeywordsKeyword struct {
	Column   string `xml:"column,omitempty"`
	Relation string `xml:"relation,omitempty"`
	Value    string `xml:"value,omitempty"`
}

// UserSort represents the <sort> element in the response.
type UserSort struct {
	Value string         `xml:",chardata"`
	Field *UserSortField `xml:"field,omitempty"`
}

// UserSortField represents the <field> element in sort.
type UserSortField struct {
	Order string `xml:"order,omitempty"`
}

// UserUsers represents the <users> meta element in the response.
type UserUsers struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// UserCount represents the <user_count> element in the response.
type UserCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
