package commands

import "encoding/xml"

// CreateNote represents a create_note command request.
type CreateNote struct {
	XMLName  xml.Name `xml:"create_note"`
	Text     string   `xml:"text"`
	NVT      *nvt     `xml:"nvt"`
	Result   *result  `xml:"result,omitempty"`
	Active   string   `xml:"active,omitempty"`
	Copy     string   `xml:"copy,omitempty"`
	Hosts    string   `xml:"hosts,omitempty"`
	Port     string   `xml:"port,omitempty"`
	Severity string   `xml:"severity,omitempty"`
	Task     *task    `xml:"task,omitempty"`
	Threat   string   `xml:"threat,omitempty"`
}

// CreateNoteResponse represents a create_note command response.
type CreateNoteResponse struct {
	XMLName    xml.Name `xml:"create_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
