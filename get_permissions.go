package gmp

import "encoding/xml"

// GetPermissionsCommand represents a get_permissions command request.
type GetPermissionsCommand struct {
	XMLName xml.Name `xml:"get_permissions"`
}

// GetPermissionsResponse represents a get_permissions command response.
type GetPermissionsResponse struct {
	XMLName     xml.Name         `xml:"get_permissions_response"`
	Status      string           `xml:"status,attr"`
	StatusText  string           `xml:"status_text,attr"`
	Permissions []PermissionInfo `xml:"permission,omitempty"`
}

type PermissionInfo struct {
	ID               string `xml:"id,attr"`
	Name             string `xml:"name"`
	Comment          string `xml:"comment,omitempty"`
	CreationTime     string `xml:"creation_time,omitempty"`
	ModificationTime string `xml:"modification_time,omitempty"`
	Writable         string `xml:"writable,omitempty"`
	InUse            string `xml:"in_use,omitempty"`
	Users            string `xml:"users,omitempty"`
}
