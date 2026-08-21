package commands

import "encoding/xml"

// Help represents a help command request.
type Help struct {
	XMLName xml.Name `xml:"help"`
	Format  string   `xml:"format,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
}

// HelpResponse represents a help command response. Text is the raw,
// undecoded payload: for the plain-text format this is the help text
// itself, for the XML format it's the opaque, unparsed <schema> subtree.
// It must use innerxml rather than chardata, since chardata silently drops
// nested elements — exactly the shape the XML format returns.
type HelpResponse struct {
	XMLName    xml.Name `xml:"help_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Text       string   `xml:",innerxml"`
}

// GetStatus returns the status and status text from the response.
func (r *HelpResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
