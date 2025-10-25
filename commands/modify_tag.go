package commands

import "encoding/xml"

// ModifyTag represents a modify_tag command request.
type ModifyTag struct {
	XMLName   xml.Name            `xml:"modify_tag"`
	TagID     string              `xml:"tag_id,attr"`
	Name      string              `xml:"name,omitempty"`
	Resources *CreateTagResources `xml:"resources,omitempty"`
	Value     string              `xml:"value,omitempty"`
	Comment   string              `xml:"comment,omitempty"`
	Active    bool                `xml:"active,omitempty"`
}

// ModifyTagResponse represents a modify_tag command response.
type ModifyTagResponse struct {
	XMLName    xml.Name `xml:"modify_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
