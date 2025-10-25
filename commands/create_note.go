package commands

import "encoding/xml"

// CreateNote is a create_note command request.
type CreateNote struct {
	XMLName  xml.Name          `xml:"create_note"`
	Text     string            `xml:"text"`
	NVT      *NoteNVT          `xml:"nvt"`
	Result   *CreateNoteResult `xml:"result,omitempty"`
	Active   int               `xml:"active,omitempty"`
	Copy     string            `xml:"copy,omitempty"`
	Hosts    string            `xml:"hosts,omitempty"`
	Port     string            `xml:"port,omitempty"`
	Severity float64           `xml:"severity,omitempty"`
	Task     *CreateNoteTask   `xml:"task,omitempty"`
	Threat   string            `xml:"threat,omitempty"`
}

// CreateNoteResponse is a create_note command response.
type CreateNoteResponse struct {
	XMLName    xml.Name `xml:"create_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// CreateNoteResult is the <result> element for create_note.
type CreateNoteResult struct {
	ID string `xml:"id,attr"`
}

// CreateNoteTask is the <task> element for create_note.
type CreateNoteTask struct {
	ID string `xml:"id,attr"`
}
