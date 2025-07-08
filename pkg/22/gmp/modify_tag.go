package gmp

import "encoding/xml"

// ModifyTagCommand represents a modify_tag command request.
type ModifyTagCommand struct {
	XMLName   xml.Name   `xml:"modify_tag"`
	TagID     string     `xml:"tag_id,attr"`
	Name      string     `xml:"name,omitempty"`
	Resources *Resources `xml:"resources,omitempty"`
	Value     string     `xml:"value,omitempty"`
	Comment   string     `xml:"comment,omitempty"`
	Active    string     `xml:"active,omitempty"`
}

// ModifyTagResponse represents a modify_tag command response.
type ModifyTagResponse struct {
	XMLName    xml.Name `xml:"modify_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
