package gmp

import "encoding/xml"

// DeleteTagCommand represents a delete_tag command request.
type DeleteTagCommand struct {
	XMLName  xml.Name `xml:"delete_tag"`
	TagID    string   `xml:"tag_id,attr"`
	Ultimate string   `xml:"ultimate,attr"`
}

// DeleteTagResponse represents a delete_tag command response.
type DeleteTagResponse struct {
	XMLName    xml.Name `xml:"delete_tag_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
