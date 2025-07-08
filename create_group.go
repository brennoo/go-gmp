package gmp

import "encoding/xml"

// CreateGroupCommand represents a create_group command request.
type CreateGroupCommand struct {
	XMLName  xml.Name  `xml:"create_group"`
	Name     string    `xml:"name"`
	Comment  string    `xml:"comment,omitempty"`
	Copy     string    `xml:"copy,omitempty"`
	Specials *Specials `xml:"specials,omitempty"`
	Users    string    `xml:"users,omitempty"`
}

type Specials struct {
	Full string `xml:"full,omitempty"`
}

// CreateGroupResponse represents a create_group command response.
type CreateGroupResponse struct {
	XMLName    xml.Name `xml:"create_group_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
