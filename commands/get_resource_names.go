package commands

import "encoding/xml"

// GetResourceNamesCommand represents a get_resource_names command request.
type GetResourceNamesCommand struct {
	XMLName    xml.Name `xml:"get_resource_names"`
	Type       string   `xml:"type,attr"`
	ResourceID string   `xml:"resource_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
}

// GetResourceNamesResponse represents a get_resource_names command response.
type GetResourceNamesResponse struct {
	XMLName    xml.Name       `xml:"get_resource_names_response"`
	Status     string         `xml:"status,attr"`
	StatusText string         `xml:"status_text,attr"`
	Type       string         `xml:"type,omitempty"`
	Resources  []ResourceName `xml:"resource"`
}

type ResourceName struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}
