package commands

import "encoding/xml"

// DeleteNote represents a delete_note command request.
type DeleteNote struct {
	XMLName  xml.Name `xml:"delete_note"`
	NoteID   string   `xml:"note_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteNoteResponse represents a delete_note command response.
type DeleteNoteResponse struct {
	XMLName    xml.Name `xml:"delete_note_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteNoteResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
