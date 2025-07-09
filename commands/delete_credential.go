package commands

import "encoding/xml"

// DeleteCredential is a request to delete a credential.
type DeleteCredential struct {
	XMLName      xml.Name `xml:"delete_credential"`
	CredentialID string   `xml:"credential_id,attr"`
	Ultimate     string   `xml:"ultimate,attr"`
}

// DeleteCredentialResponse is the response to a delete_credential command.
type DeleteCredentialResponse struct {
	XMLName    xml.Name `xml:"delete_credential_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
