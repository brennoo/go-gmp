package commands

import "encoding/xml"

// GetNotes represents a get_notes command request.
type GetNotes struct {
	XMLName xml.Name `xml:"get_notes"`
}

// GetNotesResponse represents a get_notes command response.
type GetNotesResponse struct {
	XMLName    xml.Name     `xml:"get_notes_response"`
	Status     string       `xml:"status,attr"`
	StatusText string       `xml:"status_text,attr"`
	Notes      []Note       `xml:"note,omitempty"`
	Filters    *NoteFilters `xml:"filters,omitempty"`
	Sort       *NoteSort    `xml:"sort,omitempty"`
	NoteCount  *NoteCount   `xml:"note_count,omitempty"`
}

// Note represents a <note> element in the get_notes response.
type Note struct {
	ID               string   `xml:"id,attr"`
	NVT              *NoteNVT `xml:"nvt,omitempty"`
	CreationTime     string   `xml:"creation_time,omitempty"`
	ModificationTime string   `xml:"modification_time,omitempty"`
	Writable         string   `xml:"writable,omitempty"`
	InUse            string   `xml:"in_use,omitempty"`
	Active           string   `xml:"active,omitempty"`
	Orphan           string   `xml:"orphan,omitempty"`
	Text             string   `xml:"text,omitempty"`
}

// NoteNVT represents the <nvt> element in a note.
type NoteNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name,omitempty"`
}

// NoteFilters represents the <filters> element in the response.
type NoteFilters struct {
	ID       string        `xml:"id,attr"`
	Term     string        `xml:"term"`
	Name     string        `xml:"name,omitempty"`
	Keywords *NoteKeywords `xml:"keywords,omitempty"`
}

// NoteKeywords represents the <keywords> element in filters.
type NoteKeywords struct {
	Keyword []NoteKeyword `xml:"keyword"`
}

// NoteKeyword represents a <keyword> element in filters keywords.
type NoteKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// NoteSort represents the <sort> element in the response.
type NoteSort struct {
	Value string         `xml:",chardata"`
	Field *NoteSortField `xml:"field,omitempty"`
}

// NoteSortField represents the <field> element in sort.
type NoteSortField struct {
	Order string `xml:"order"`
}

// NoteCount represents the <note_count> element in the response.
type NoteCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
