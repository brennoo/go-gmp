package gmp

import "encoding/xml"

// ModifyUserCommand represents a modify_user command request.
type ModifyUserCommand struct {
	XMLName  xml.Name            `xml:"modify_user"`
	UserID   string              `xml:"user_id,attr"`
	Name     string              `xml:"name"`
	NewName  string              `xml:"new_name,omitempty"`
	Comment  string              `xml:"comment,omitempty"`
	Password *ModifyUserPassword `xml:"password,omitempty"`
	Roles    []CreateUserRole    `xml:"role,omitempty"`
	Hosts    *CreateUserHosts    `xml:"hosts,omitempty"`
	Sources  *CreateUserSources  `xml:"sources,omitempty"`
}

type ModifyUserPassword struct {
	Modify string `xml:"modify,attr"`
	Text   string `xml:",chardata"`
}

// ModifyUserResponse represents a modify_user command response.
type ModifyUserResponse struct {
	XMLName    xml.Name `xml:"modify_user_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
