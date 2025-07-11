package commands

import "encoding/xml"

// CreateCredential represents a GMP create_credential command request.
type CreateCredential struct {
	XMLName       xml.Name                 `xml:"create_credential"`
	Name          string                   `xml:"name"`
	Comment       string                   `xml:"comment,omitempty"`
	Copy          string                   `xml:"copy,omitempty"`
	AllowInsecure string                   `xml:"allow_insecure,omitempty"`
	Certificate   string                   `xml:"certificate,omitempty"`
	KDC           string                   `xml:"kdc,omitempty"`
	KDCs          *CreateCredentialKDCs    `xml:"kdcs,omitempty"`
	Key           *CreateCredentialKey     `xml:"key,omitempty"`
	Login         string                   `xml:"login,omitempty"`
	Password      string                   `xml:"password,omitempty"`
	AuthAlgorithm string                   `xml:"auth_algorithm,omitempty"`
	Community     string                   `xml:"community,omitempty"`
	Privacy       *CreateCredentialPrivacy `xml:"privacy,omitempty"`
	Realm         string                   `xml:"realm,omitempty"`
	Type          string                   `xml:"type,omitempty"`
}

// CreateCredentialResponse represents a GMP create_credential command response.
type CreateCredentialResponse struct {
	XMLName    xml.Name `xml:"create_credential_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// CreateCredentialKDCs represents a list of Kerberos Key Distribution Centers (KDCs).
type CreateCredentialKDCs struct {
	KDC []string `xml:"kdc"`
}

// CredentialKey represents the key element for key-based credentials.
type CreateCredentialKey struct {
	Phrase  string `xml:"phrase,omitempty"`
	Private string `xml:"private,omitempty"`
	Public  string `xml:"public,omitempty"`
}

// CreateCredentialPrivacy represents SNMP privacy settings.
type CreateCredentialPrivacy struct {
	Algorithm string `xml:"algorithm,omitempty"`
	Password  string `xml:"password,omitempty"`
}
