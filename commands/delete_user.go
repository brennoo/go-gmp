package commands

import "encoding/xml"

// DeleteUser represents a delete_user command request.
type DeleteUser struct {
	XMLName       xml.Name `xml:"delete_user"`
	UserID        string   `xml:"user_id,attr,omitempty"`
	Name          string   `xml:"name,attr,omitempty"`
	InheritorID   string   `xml:"inheritor_id,attr,omitempty"`
	InheritorName string   `xml:"inheritor_name,attr,omitempty"`
}

// DeleteUserResponse represents a delete_user command response.
type DeleteUserResponse struct {
	XMLName    xml.Name `xml:"delete_user_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteUserResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
