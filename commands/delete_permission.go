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

// GetStatus returns the status and status text from the response.
func (r *DeletePermissionResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
