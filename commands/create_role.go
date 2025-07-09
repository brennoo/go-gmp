package commands

import "encoding/xml"

// CreateRole represents a create_role command request.
type CreateRole struct {
	XMLName xml.Name `xml:"create_role"`
	Name    string   `xml:"name"`
	Comment string   `xml:"comment,omitempty"`
}

// CreateRoleResponse represents a create_role command response.
type CreateRoleResponse struct {
	XMLName    xml.Name `xml:"create_role_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr,omitempty"`
}
