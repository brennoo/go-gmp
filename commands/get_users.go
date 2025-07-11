package commands

import "encoding/xml"

// GetUsers represents a get_users command request.
type GetUsers struct {
	XMLName xml.Name `xml:"get_users"`
}

// GetUsersResponse represents a get_users command response.
type GetUsersResponse struct {
	XMLName    xml.Name `xml:"get_users_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Users      []User   `xml:"user,omitempty"`
}

// User represents a <user> element in the get_users response.
type User struct {
	ID               string           `xml:"id,attr,omitempty"`
	Name             string           `xml:"name"`
	Comment          string           `xml:"comment,omitempty"`
	CreationTime     string           `xml:"creation_time,omitempty"`
	ModificationTime string           `xml:"modification_time,omitempty"`
	Writable         string           `xml:"writable,omitempty"`
	InUse            string           `xml:"in_use,omitempty"`
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

// UserSources represents the <sources> element in a user.
type UserSources struct {
	Source []string `xml:"source"`
}
