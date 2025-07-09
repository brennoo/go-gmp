package commands

import "encoding/xml"

// CreateTagCommand represents a create_tag command request.
type CreateTagCommand struct {
	XMLName   xml.Name   `xml:"create_tag"`
	Name      string     `xml:"name"`
	Resources *Resources `xml:"resources,omitempty"`
	Copy      string     `xml:"copy,omitempty"`
	Value     string     `xml:"value,omitempty"`
	Comment   string     `xml:"comment,omitempty"`
	Active    string     `xml:"active,omitempty"`
}

type Resources struct {
	Resource *Resource `xml:"resource,omitempty"`
	Type     string    `xml:"type"`
}

type Resource struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTagResponse represents a create_tag command response.
type CreateTagResponse struct {
	XMLName    xml.Name `xml:"create_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
