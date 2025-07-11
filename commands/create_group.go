package commands

import "encoding/xml"

// CreateGroup is a create_group command request.
type CreateGroup struct {
	XMLName  xml.Name  `xml:"create_group"`
	Name     string    `xml:"name"`
	Comment  string    `xml:"comment,omitempty"`
	Copy     string    `xml:"copy,omitempty"`
	Specials *Specials `xml:"specials,omitempty"`
	Users    string    `xml:"users,omitempty"`
}

// Specials is the <specials> element for create_group.
type Specials struct {
	Full string `xml:"full,omitempty"`
}

// CreateGroupResponse is a create_group command response.
type CreateGroupResponse struct {
	XMLName    xml.Name `xml:"create_group_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
