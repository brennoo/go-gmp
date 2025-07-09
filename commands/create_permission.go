package commands

import "encoding/xml"

// CreatePermission represents a create_permission command request.
type CreatePermission struct {
	XMLName  xml.Name            `xml:"create_permission"`
	Name     string              `xml:"name"`
	Subject  *permissionSubject  `xml:"subject"`
	Resource *permissionResource `xml:"resource,omitempty"`
	Copy     string              `xml:"copy,omitempty"`
	Comment  string              `xml:"comment,omitempty"`
}

type permissionSubject struct {
	ID   string `xml:"id,attr,omitempty"`
	Type string `xml:"type"`
}

func NewPermissionSubject(id, typ string) *permissionSubject {
	return &permissionSubject{ID: id, Type: typ}
}

type permissionResource struct {
	ID   string `xml:"id,attr,omitempty"`
	Type string `xml:"type,omitempty"`
}

func NewPermissionResource(id, typ string) *permissionResource {
	return &permissionResource{ID: id, Type: typ}
}

// CreatePermissionResponse represents a create_permission command response.
type CreatePermissionResponse struct {
	XMLName    xml.Name `xml:"create_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr,omitempty"`
}
