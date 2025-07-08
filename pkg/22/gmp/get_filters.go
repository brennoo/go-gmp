package gmp

import "encoding/xml"

// GetFiltersCommand represents a get_filters command request.
type GetFiltersCommand struct {
	XMLName  xml.Name `xml:"get_filters"`
	FilterID string   `xml:"filter_id,attr,omitempty"`
}

// GetFiltersResponse represents a get_filters command response.
type GetFiltersResponse struct {
	XMLName    xml.Name      `xml:"get_filters_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Filters    []FilterEntry `xml:"filter"`
}

// FilterEntry represents a filter entry in the get_filters response.
type FilterEntry struct {
	ID               string `xml:"id,attr"`
	Name             string `xml:"name"`
	Comment          string `xml:"comment"`
	Term             string `xml:"term"`
	Type             string `xml:"type"`
	InUse            string `xml:"in_use"`
	Writable         string `xml:"writable"`
	CreationTime     string `xml:"creation_time"`
	ModificationTime string `xml:"modification_time"`
}
