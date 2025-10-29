package commands

import "encoding/xml"

// CreateRole represents a create_role command request.
type CreateRole struct {
	XMLName xml.Name `xml:"create_role"`
	Name    string   `xml:"name"`
	Comment string   `xml:"comment,omitempty"`
	Copy    string   `xml:"copy,omitempty"`
	Users   string   `xml:"users,omitempty"`
}

// CreateRoleResponse represents a create_role command response.
type CreateRoleResponse struct {
	XMLName    xml.Name `xml:"create_role_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateRoleResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
