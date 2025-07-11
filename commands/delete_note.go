package commands

import "encoding/xml"

// DeleteNote represents a delete_note command request.
type DeleteNote struct {
	XMLName  xml.Name `xml:"delete_note"`
	NoteID   string   `xml:"note_id,attr"`
	Ultimate string   `xml:"ultimate,attr"`
}

// DeleteNoteResponse represents a delete_note command response.
type DeleteNoteResponse struct {
	XMLName    xml.Name `xml:"delete_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
