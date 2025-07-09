package commands

import "encoding/xml"

// ModifyPermissionCommand represents a modify_permission command request.
type ModifyPermissionCommand struct {
	XMLName      xml.Name                  `xml:"modify_permission"`
	PermissionID string                    `xml:"permission_id,attr"`
	Name         string                    `xml:"name,omitempty"`
	Comment      string                    `xml:"comment,omitempty"`
	Resource     *CreatePermissionResource `xml:"resource,omitempty"`
	Subject      *CreatePermissionSubject  `xml:"subject,omitempty"`
}

// ModifyPermissionResponse represents a modify_permission command response.
type ModifyPermissionResponse struct {
	XMLName    xml.Name `xml:"modify_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
