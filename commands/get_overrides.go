package commands

import (
	"encoding/xml"
)

// GetOverridesCommand represents a get_overrides command request.
type GetOverridesCommand struct {
	XMLName    xml.Name `xml:"get_overrides"`
	OverrideID string   `xml:"override_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
	FiltID     string   `xml:"filt_id,attr,omitempty"`
	Trash      string   `xml:"trash,attr,omitempty"`
	Details    string   `xml:"details,attr,omitempty"`
	Result     string   `xml:"result,attr,omitempty"`
}

// OverrideFilters represents filter information in the response.
type OverrideFilters struct {
	ID       string `xml:"id,attr"`
	Term     string `xml:"term"`
	Name     string `xml:"name,omitempty"`
	Keywords string `xml:"keywords,omitempty"`
}

// OverrideSort represents sort information in the response.
type OverrideSort struct {
	Field string `xml:"field"`
	Order string `xml:"order"`
}

// OverridePagination represents pagination information in the response.
type OverridePagination struct {
	Start string `xml:"start,attr"`
	Max   string `xml:"max,attr"`
}

// OverrideCount represents count information in the response.
type OverrideCount struct {
	Filtered string `xml:"filtered"`
	Page     string `xml:"page"`
}

// GetOverridesResponse represents a get_overrides command response.
type GetOverridesResponse struct {
	XMLName    xml.Name            `xml:"get_overrides_response"`
	Status     string              `xml:"status,attr"`
	StatusText string              `xml:"status_text,attr"`
	Overrides  []Override          `xml:"override,omitempty"`
	Filters    *OverrideFilters    `xml:"filters,omitempty"`
	Sort       *OverrideSort       `xml:"sort,omitempty"`
	Pagination *OverridePagination `xml:"overrides,omitempty"`
	Count      *OverrideCount      `xml:"override_count,omitempty"`
}
