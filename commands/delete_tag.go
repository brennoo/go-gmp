package commands

import "encoding/xml"

// DeleteTag represents a delete_tag command request.
type DeleteTag struct {
	XMLName  xml.Name `xml:"delete_tag"`
	TagID    string   `xml:"tag_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteTagResponse represents a delete_tag command response.
type DeleteTagResponse struct {
	XMLName    xml.Name `xml:"delete_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteTagResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
