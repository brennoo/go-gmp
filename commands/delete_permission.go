package commands

import "encoding/xml"

// DeletePermission represents a delete_permission command request.
type DeletePermission struct {
	XMLName      xml.Name `xml:"delete_permission"`
	PermissionID string   `xml:"permission_id,attr"`
	Ultimate     bool     `xml:"ultimate,attr"`
}

// DeletePermissionResponse represents a delete_permission command response.
type DeletePermissionResponse struct {
	XMLName    xml.Name `xml:"delete_permission_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
