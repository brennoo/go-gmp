package commands

import "encoding/xml"

// DeleteFilter represents a delete_filter command request.
type DeleteFilter struct {
	XMLName  xml.Name `xml:"delete_filter"`
	FilterID string   `xml:"filter_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteFilterResponse represents a delete_filter command response.
type DeleteFilterResponse struct {
	XMLName    xml.Name `xml:"delete_filter_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteFilterResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
