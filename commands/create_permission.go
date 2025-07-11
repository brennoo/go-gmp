package commands

import "encoding/xml"

// CreatePermission is a create_permission command request.
type CreatePermission struct {
	XMLName  xml.Name            `xml:"create_permission"`
	Name     string              `xml:"name"`
	Subject  *PermissionSubject  `xml:"subject"`
	Resource *PermissionResource `xml:"resource,omitempty"`
	Copy     string              `xml:"copy,omitempty"`
	Comment  string              `xml:"comment,omitempty"`
}

// CreatePermissionResponse is a create_permission command response.
type CreatePermissionResponse struct {
	XMLName    xml.Name `xml:"create_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr,omitempty"`
}
