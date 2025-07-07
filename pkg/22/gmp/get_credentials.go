package gmp

import "encoding/xml"

// GetCredentialsCommand represents a get_credentials command request to GMP.
// It can be used to fetch one or more credentials, optionally filtered by credential_id, format, or targets.
type GetCredentialsCommand struct {
	XMLName      xml.Name `xml:"get_credentials"`
	CredentialID string   `xml:"credential_id,attr,omitempty"`
	Format       string   `xml:"format,attr,omitempty"`
	Targets      string   `xml:"targets,attr,omitempty"`
}

// GetCredentialsResponse represents a get_credentials command response from GMP.
// It contains a list of credentials and the response status.
type GetCredentialsResponse struct {
	XMLName     xml.Name            `xml:"get_credentials_response"`
	Status      string              `xml:"status,attr"`
	StatusText  string              `xml:"status_text,attr"`
	Credentials []CredentialWrapper `xml:"credential"`
	Truncated   string              `xml:"truncated,omitempty"`
}

// CredentialWrapper represents a <credential> element in the get_credentials response.
type CredentialWrapper struct {
	ID               string             `xml:"id,attr"`
	Name             string             `xml:"name,omitempty"`
	Login            string             `xml:"login,omitempty"`
	Comment          string             `xml:"comment,omitempty"`
	CreationTime     string             `xml:"creation_time,omitempty"`
	ModificationTime string             `xml:"modification_time,omitempty"`
	Writable         string             `xml:"writable,omitempty"`
	InUse            string             `xml:"in_use,omitempty"`
	Type             string             `xml:"type,omitempty"`
	FullType         string             `xml:"full_type,omitempty"`
	Formats          *CredentialFormats `xml:"formats,omitempty"`
	Targets          *CredentialTargets `xml:"targets,omitempty"`
	PublicKey        string             `xml:"public_key,omitempty"`
	Package          *CredentialPackage `xml:"package,omitempty"`
}

// CredentialFormats represents the <formats> element containing multiple <format>.
type CredentialFormats struct {
	Formats []string `xml:"format"`
}

// CredentialTargets represents the <targets> element containing multiple <target>.
type CredentialTargets struct {
	Targets []CredentialTarget `xml:"target"`
}

// CredentialTarget represents a <target> element under <targets>.
type CredentialTarget struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,omitempty"`
}

// CredentialPackage represents the <package> element under <credential>.
type CredentialPackage struct {
	Format string `xml:"format,attr,omitempty"`
	Value  string `xml:",chardata"`
}
