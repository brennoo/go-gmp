package gmp

import "encoding/xml"

// ModifyFilterCommand represents a modify_filter command request.
type ModifyFilterCommand struct {
	XMLName  xml.Name `xml:"modify_filter"`
	FilterID string   `xml:"filter_id,attr"`
	Comment  string   `xml:"comment,omitempty"`
	Name     string   `xml:"name,omitempty"`
	Term     string   `xml:"term,omitempty"`
	Type     string   `xml:"type,omitempty"`
}

// ModifyFilterResponse represents a modify_filter command response.
type ModifyFilterResponse struct {
	XMLName    xml.Name `xml:"modify_filter_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
