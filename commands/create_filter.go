package commands

import "encoding/xml"

// CreateFilter represents a create_filter command request.
type CreateFilter struct {
	XMLName xml.Name `xml:"create_filter"`
	Name    string   `xml:"name"`
	Comment string   `xml:"comment,omitempty"`
	Copy    string   `xml:"copy,omitempty"`
	Term    string   `xml:"term,omitempty"`
	Type    string   `xml:"type,omitempty"`
}

// CreateFilterResponse represents a create_filter command response.
type CreateFilterResponse struct {
	XMLName    xml.Name `xml:"create_filter_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
