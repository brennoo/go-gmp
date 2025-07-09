package commands

import "encoding/xml"

// GetGroups represents a get_groups command request.
type GetGroups struct {
	XMLName xml.Name `xml:"get_groups"`
}

// GetGroupsResponse represents a get_groups command response.
type GetGroupsResponse struct {
	XMLName    xml.Name    `xml:"get_groups_response"`
	Status     string      `xml:"status,attr"`
	StatusText string      `xml:"status_text,attr"`
	Groups     []GroupInfo `xml:"group,omitempty"`
}

type GroupInfo struct {
	ID               string `xml:"id,attr"`
	Name             string `xml:"name"`
	Comment          string `xml:"comment,omitempty"`
	CreationTime     string `xml:"creation_time,omitempty"`
	ModificationTime string `xml:"modification_time,omitempty"`
	Writable         string `xml:"writable,omitempty"`
	InUse            string `xml:"in_use,omitempty"`
	Users            string `xml:"users,omitempty"`
}
