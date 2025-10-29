package commands

import "encoding/xml"

// ModifyPermission is a modify_permission command request.
type ModifyPermission struct {
	XMLName      xml.Name            `xml:"modify_permission"`
	PermissionID string              `xml:"permission_id,attr"`
	Name         string              `xml:"name,omitempty"`
	Comment      string              `xml:"comment,omitempty"`
	Resource     *PermissionResource `xml:"resource,omitempty"`
	Subject      *PermissionSubject  `xml:"subject,omitempty"`
}

// ModifyPermissionResponse is a modify_permission command response.
type ModifyPermissionResponse struct {
	XMLName    xml.Name `xml:"modify_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyPermissionResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
