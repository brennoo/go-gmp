package commands

import "encoding/xml"

// ModifyPermission represents a modify_permission command request.
type ModifyPermission struct {
	XMLName      xml.Name            `xml:"modify_permission"`
	PermissionID string              `xml:"permission_id,attr"`
	Name         string              `xml:"name,omitempty"`
	Comment      string              `xml:"comment,omitempty"`
	Resource     *permissionResource `xml:"resource,omitempty"`
	Subject      *permissionSubject  `xml:"subject,omitempty"`
}

// ModifyPermissionResponse represents a modify_permission command response.
type ModifyPermissionResponse struct {
	XMLName    xml.Name `xml:"modify_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
