package commands

import "encoding/xml"

// GetRoles represents a get_roles command request.
type GetRoles struct {
	XMLName xml.Name `xml:"get_roles"`
}

// GetRolesResponse represents a get_roles command response.
type GetRolesResponse struct {
	XMLName    xml.Name `xml:"get_roles_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Roles      []Role   `xml:"role,omitempty"`
}

type Role struct {
	ID               string `xml:"id,attr"`
	Name             string `xml:"name"`
	Comment          string `xml:"comment,omitempty"`
	CreationTime     string `xml:"creation_time,omitempty"`
	ModificationTime string `xml:"modification_time,omitempty"`
	Writable         string `xml:"writable,omitempty"`
	InUse            string `xml:"in_use,omitempty"`
	Users            string `xml:"users,omitempty"`
}
