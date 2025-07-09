package commands

import "encoding/xml"

// GetNotes represents a get_notes command request.
type GetNotes struct {
	XMLName xml.Name `xml:"get_notes"`
}

// GetNotesResponse represents a get_notes command response.
type GetNotesResponse struct {
	XMLName    xml.Name   `xml:"get_notes_response"`
	Status     string     `xml:"status,attr"`
	StatusText string     `xml:"status_text,attr"`
	Notes      []NoteInfo `xml:"note,omitempty"`
}

type NoteInfo struct {
	ID               string      `xml:"id,attr"`
	NVT              *GetNoteNVT `xml:"nvt,omitempty"`
	CreationTime     string      `xml:"creation_time,omitempty"`
	ModificationTime string      `xml:"modification_time,omitempty"`
	Writable         string      `xml:"writable,omitempty"`
	InUse            string      `xml:"in_use,omitempty"`
	Active           string      `xml:"active,omitempty"`
	Orphan           string      `xml:"orphan,omitempty"`
	Text             string      `xml:"text,omitempty"`
}

type GetNoteNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name,omitempty"`
}
