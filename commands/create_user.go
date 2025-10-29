package commands

import "encoding/xml"

// CreateUser represents a create_user command request.
type CreateUser struct {
	XMLName  xml.Name     `xml:"create_user"`
	Name     string       `xml:"name"`
	Copy     string       `xml:"copy,omitempty"`
	Comment  string       `xml:"comment,omitempty"`
	Hosts    *UserHosts   `xml:"hosts,omitempty"`
	Password string       `xml:"password,omitempty"`
	Roles    []*UserRole  `xml:"role,omitempty"`
	Groups   *UserGroups  `xml:"groups,omitempty"`
	Sources  *UserSources `xml:"sources,omitempty"`
}

// CreateUserResponse represents a create_user command response.
type CreateUserResponse struct {
	XMLName    xml.Name `xml:"create_user_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateUserResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
