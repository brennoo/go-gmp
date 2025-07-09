package commands

import "encoding/xml"

// CreatePermission represents a create_permission command request.
type CreatePermission struct {
	XMLName  xml.Name                  `xml:"create_permission"`
	Name     string                    `xml:"name"`
	Subject  *CreatePermissionSubject  `xml:"subject"`
	Resource *CreatePermissionResource `xml:"resource,omitempty"`
	Copy     string                    `xml:"copy,omitempty"`
	Comment  string                    `xml:"comment,omitempty"`
}

type CreatePermissionSubject struct {
	ID   string `xml:"id,attr,omitempty"`
	Type string `xml:"type"`
}

type CreatePermissionResource struct {
	ID   string `xml:"id,attr,omitempty"`
	Type string `xml:"type,omitempty"`
}

// CreatePermissionResponse represents a create_permission command response.
type CreatePermissionResponse struct {
	XMLName    xml.Name `xml:"create_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr,omitempty"`
}
