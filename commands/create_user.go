package commands

import "encoding/xml"

// CreateUser represents a create_user command request.
type CreateUser struct {
	XMLName  xml.Name           `xml:"create_user"`
	Name     string             `xml:"name"`
	Copy     string             `xml:"copy,omitempty"`
	Comment  string             `xml:"comment,omitempty"`
	Hosts    *CreateUserHosts   `xml:"hosts,omitempty"`
	Password string             `xml:"password,omitempty"`
	Roles    []CreateUserRole   `xml:"role,omitempty"`
	Groups   *CreateUserGroups  `xml:"groups,omitempty"`
	Sources  *CreateUserSources `xml:"sources,omitempty"`
}

// CreateUserHosts represents the hosts element for create_user.
type CreateUserHosts struct {
	Allow string `xml:"allow,attr,omitempty"`
	Text  string `xml:",chardata"`
}

// CreateUserRole represents a role element for create_user.
type CreateUserRole struct {
	ID string `xml:"id,attr"`
}

// CreateUserGroups represents the groups element for create_user.
type CreateUserGroups struct {
	Groups []CreateUserGroup `xml:"group,omitempty"`
}

type CreateUserGroup struct {
	ID string `xml:"id,attr"`
}

// CreateUserSources represents the sources element for create_user.
type CreateUserSources struct {
	Source *CreateUserSource `xml:"source,omitempty"`
}

type CreateUserSource struct {
	Type string `xml:",chardata"`
}

// CreateUserResponse represents a create_user command response.
type CreateUserResponse struct {
	XMLName    xml.Name `xml:"create_user_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
