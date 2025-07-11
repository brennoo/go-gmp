package commands

import "encoding/xml"

// ModifyUser represents a modify_user command request.
type ModifyUser struct {
	XMLName  xml.Name      `xml:"modify_user"`
	UserID   string        `xml:"user_id,attr"`
	Name     string        `xml:"name"`
	NewName  string        `xml:"new_name,omitempty"`
	Comment  string        `xml:"comment,omitempty"`
	Password *UserPassword `xml:"password,omitempty"`
	Roles    []*UserRole   `xml:"role,omitempty"`
	Hosts    *UserHosts    `xml:"hosts,omitempty"`
	Sources  *UserSources  `xml:"sources,omitempty"`
}

// ModifyUserResponse represents a modify_user command response.
type ModifyUserResponse struct {
	XMLName    xml.Name `xml:"modify_user_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

type UserPassword struct {
	Modify string `xml:"modify,attr,omitempty"`
	Text   string `xml:",chardata"`
}
