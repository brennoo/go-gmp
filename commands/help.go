package commands

import "encoding/xml"

// Help represents a help command request.
type Help struct {
	XMLName xml.Name `xml:"help"`
	Format  string   `xml:"format,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
}

// HelpResponse represents a help command response.
type HelpResponse struct {
	XMLName    xml.Name `xml:"help_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Text       string   `xml:",chardata"`
}
