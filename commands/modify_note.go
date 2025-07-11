package commands

import "encoding/xml"

// ModifyNote represents a modify_note command request.
type ModifyNote struct {
	XMLName  xml.Name          `xml:"modify_note"`
	NoteID   string            `xml:"note_id,attr"`
	Text     string            `xml:"text,omitempty"`
	NVT      *NoteNVT          `xml:"nvt,omitempty"`
	Result   *ModifyNoteResult `xml:"result,omitempty"`
	Active   string            `xml:"active,omitempty"`
	Hosts    string            `xml:"hosts,omitempty"`
	Port     string            `xml:"port,omitempty"`
	Severity string            `xml:"severity,omitempty"`
	Task     *ModifyNoteTask   `xml:"task,omitempty"`
}

// ModifyNoteResponse represents a modify_note command response.
type ModifyNoteResponse struct {
	XMLName    xml.Name `xml:"modify_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// ModifyNoteResult represents the <result> element for modify_note.
type ModifyNoteResult struct {
	ID string `xml:"id,attr"`
}

// ModifyNoteTask represents the <task> element for modify_note.
type ModifyNoteTask struct {
	ID string `xml:"id,attr"`
}
