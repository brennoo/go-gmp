package commands

import "encoding/xml"

// ModifyGroup represents a modify_group command request.
type ModifyGroup struct {
	XMLName xml.Name `xml:"modify_group"`
	GroupID string   `xml:"group_id,attr"`
	Name    string   `xml:"name,omitempty"`
	Comment string   `xml:"comment,omitempty"`
	Users   string   `xml:"users,omitempty"`
}

// ModifyGroupResponse represents a modify_group command response.
type ModifyGroupResponse struct {
	XMLName    xml.Name `xml:"modify_group_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
