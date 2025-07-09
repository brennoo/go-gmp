package commands

import "encoding/xml"

// GetNvtFamiliesCommand represents a get_nvt_families command request.
type GetNvtFamiliesCommand struct {
	XMLName   xml.Name `xml:"get_nvt_families"`
	SortOrder string   `xml:"sort_order,attr,omitempty"`
}

// GetNvtFamiliesResponse represents a get_nvt_families command response.
type GetNvtFamiliesResponse struct {
	XMLName    xml.Name `xml:"get_nvt_families_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Families   []Family `xml:"families>family,omitempty"`
}

// Family represents an NVT family.
type Family struct {
	Name        string `xml:"name,omitempty"`
	MaxNvtCount int    `xml:"max_nvt_count,omitempty"`
}
