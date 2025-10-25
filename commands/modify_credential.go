package commands

import "encoding/xml"

// ModifyCredential represents a modify_credential command request to GMP.
type ModifyCredential struct {
	XMLName       xml.Name                 `xml:"modify_credential"`
	CredentialID  string                   `xml:"credential_id,attr"`
	Comment       string                   `xml:"comment,omitempty"`
	Name          string                   `xml:"name,omitempty"`
	AllowInsecure bool                     `xml:"allow_insecure,omitempty"`
	Certificate   string                   `xml:"certificate,omitempty"`
	KDC           string                   `xml:"kdc,omitempty"`
	KDCs          *CreateCredentialKDCs    `xml:"kdcs,omitempty"`
	Key           *CreateCredentialKey     `xml:"key,omitempty"`
	Login         string                   `xml:"login,omitempty"`
	Password      string                   `xml:"password,omitempty"`
	Community     string                   `xml:"community,omitempty"`
	AuthAlgorithm string                   `xml:"auth_algorithm,omitempty"`
	Privacy       *CreateCredentialPrivacy `xml:"privacy,omitempty"`
	Realm         string                   `xml:"realm,omitempty"`
}

// ModifyCredentialResponse represents a modify_credential command response from GMP.
type ModifyCredentialResponse struct {
	XMLName    xml.Name `xml:"modify_credential_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
