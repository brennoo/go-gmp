package commands

import "encoding/xml"

// EmptyTrashcan represents an empty_trashcan command request.
type EmptyTrashcan struct {
	XMLName xml.Name `xml:"empty_trashcan"`
}

// EmptyTrashcanResponse represents an empty_trashcan command response.
type EmptyTrashcanResponse struct {
	XMLName    xml.Name `xml:"empty_trashcan_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *EmptyTrashcanResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
