package commands

import "encoding/xml"

// DeleteCredential represents a delete_credential command request.
type DeleteCredential struct {
	XMLName      xml.Name `xml:"delete_credential"`
	CredentialID string   `xml:"credential_id,attr"`
	Ultimate     bool     `xml:"ultimate,attr"`
}

// DeleteCredentialResponse represents a delete_credential command response.
type DeleteCredentialResponse struct {
	XMLName    xml.Name `xml:"delete_credential_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteCredentialResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
