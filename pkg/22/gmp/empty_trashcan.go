package gmp

import "encoding/xml"

// EmptyTrashcanCommand represents an empty_trashcan command request.
type EmptyTrashcanCommand struct {
	XMLName xml.Name `xml:"empty_trashcan"`
}

// EmptyTrashcanResponse represents an empty_trashcan command response.
type EmptyTrashcanResponse struct {
	XMLName    xml.Name `xml:"empty_trashcan_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
