package commands

import "encoding/xml"

// CreateNoteCommand represents a create_note command request.
type CreateNoteCommand struct {
	XMLName  xml.Name          `xml:"create_note"`
	Text     string            `xml:"text"`
	NVT      *CreateNoteNVT    `xml:"nvt"`
	Result   *CreateNoteResult `xml:"result,omitempty"`
	Active   string            `xml:"active,omitempty"`
	Copy     string            `xml:"copy,omitempty"`
	Hosts    string            `xml:"hosts,omitempty"`
	Port     string            `xml:"port,omitempty"`
	Severity string            `xml:"severity,omitempty"`
	Task     *CreateNoteTask   `xml:"task,omitempty"`
	Threat   string            `xml:"threat,omitempty"`
}

type CreateNoteNVT struct {
	OID string `xml:"oid,attr"`
}

type CreateNoteResult struct {
	ID string `xml:"id,attr"`
}

type CreateNoteTask struct {
	ID string `xml:"id,attr"`
}

// CreateNoteResponse represents a create_note command response.
type CreateNoteResponse struct {
	XMLName    xml.Name `xml:"create_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
