package commands

import "encoding/xml"

// ModifyRole represents a modify_role command request.
type ModifyRole struct {
	XMLName xml.Name `xml:"modify_role"`
	RoleID  string   `xml:"role_id,attr"`
	Name    string   `xml:"name,omitempty"`
	Comment string   `xml:"comment,omitempty"`
	Users   string   `xml:"users,omitempty"`
}

// ModifyRoleResponse represents a modify_role command response.
type ModifyRoleResponse struct {
	XMLName    xml.Name `xml:"modify_role_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyRoleResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
