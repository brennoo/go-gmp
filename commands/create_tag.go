package commands

import "encoding/xml"

// CreateTag represents a create_tag command request.
type CreateTag struct {
	XMLName   xml.Name            `xml:"create_tag"`
	Name      string              `xml:"name"`
	Resources *CreateTagResources `xml:"resources,omitempty"`
	Copy      string              `xml:"copy,omitempty"`
	Value     string              `xml:"value,omitempty"`
	Comment   string              `xml:"comment,omitempty"`
	Active    bool                `xml:"active,omitempty"`
}

// CreateTagResponse represents a create_tag command response.
type CreateTagResponse struct {
	XMLName    xml.Name `xml:"create_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// CreateTagResources represents the <resources> element for create_tag.
type CreateTagResources struct {
	Filter    string        `xml:"filter,attr,omitempty"`
	Resources []TagResource `xml:"resource,omitempty"`
	Type      string        `xml:"type"`
}
