package gmp

import "encoding/xml"

// ModifyNoteCommand represents a modify_note command request.
type ModifyNoteCommand struct {
	XMLName  xml.Name          `xml:"modify_note"`
	NoteID   string            `xml:"note_id,attr"`
	Text     string            `xml:"text,omitempty"`
	NVT      *ModifyNoteNVT    `xml:"nvt,omitempty"`
	Result   *ModifyNoteResult `xml:"result,omitempty"`
	Active   string            `xml:"active,omitempty"`
	Hosts    string            `xml:"hosts,omitempty"`
	Port     string            `xml:"port,omitempty"`
	Severity string            `xml:"severity,omitempty"`
	Task     *ModifyNoteTask   `xml:"task,omitempty"`
}

type ModifyNoteNVT struct {
	OID string `xml:"oid,attr"`
}

type ModifyNoteResult struct {
	ID string `xml:"id,attr"`
}

type ModifyNoteTask struct {
	ID string `xml:"id,attr"`
}

// ModifyNoteResponse represents a modify_note command response.
type ModifyNoteResponse struct {
	XMLName    xml.Name `xml:"modify_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
