package commands

import "encoding/xml"

// GetNVTFamilies represents a get_nvt_families command request.
type GetNVTFamilies struct {
	XMLName   xml.Name `xml:"get_nvt_families"`
	SortOrder string   `xml:"sort_order,attr,omitempty"`
}

// GetNVTFamiliesResponse represents a get_nvt_families command response.
type GetNVTFamiliesResponse struct {
	XMLName    xml.Name    `xml:"get_nvt_families_response"`
	Status     string      `xml:"status,attr"`
	StatusText string      `xml:"status_text,attr"`
	Families   []NVTFamily `xml:"families>family,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetNVTFamiliesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// NVTFamily represents an NVT family in the response.
type NVTFamily struct {
	Name        string `xml:"name"`
	MaxNVTCount int    `xml:"max_nvt_count"`
}
