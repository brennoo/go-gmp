package gmp

import "encoding/xml"

// GetUsersCommand represents a get_users command request.
type GetUsersCommand struct {
	XMLName xml.Name `xml:"get_users"`
}

// GetUsersResponse represents a get_users command response.
type GetUsersResponse struct {
	XMLName    xml.Name      `xml:"get_users_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Users      []UserWrapper `xml:"user,omitempty"`
}

type UserWrapper struct {
	Name    string             `xml:"name"`
	Role    *UserRole          `xml:"role,omitempty"`
	Hosts   *CreateUserHosts   `xml:"hosts,omitempty"`
	Sources *CreateUserSources `xml:"sources,omitempty"`
}

type UserRole struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}
