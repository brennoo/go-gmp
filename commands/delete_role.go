package commands

import "encoding/xml"

// DeleteRole represents a delete_role command request.
type DeleteRole struct {
	XMLName  xml.Name `xml:"delete_role"`
	RoleID   string   `xml:"role_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteRoleResponse represents a delete_role command response.
type DeleteRoleResponse struct {
	XMLName    xml.Name `xml:"delete_role_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteRoleResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
