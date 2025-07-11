package commands

import "encoding/xml"

// GetTLSCertificates represents a get_tls_certificates command request.
type GetTLSCertificates struct {
	XMLName                xml.Name `xml:"get_tls_certificates"`
	IncludeCertificateData string   `xml:"include_certificate_data,attr,omitempty"` // protocol: boolean as string
}

// GetTLSCertificatesResponse represents a get_tls_certificates command response.
type GetTLSCertificatesResponse struct {
	XMLName         xml.Name         `xml:"get_tls_certificates_response"`
	Status          string           `xml:"status,attr"`
	StatusText      string           `xml:"status_text,attr"`
	TLSCertificates []TLSCertificate `xml:"tls_certificate"`
}

// TLSCertificate represents a <tls_certificate> element in the get_tls_certificates response.
type TLSCertificate struct {
	ID                string                     `xml:"id,attr"`
	Owner             *TLSCertificateOwner       `xml:"owner,omitempty"`
	Name              string                     `xml:"name,omitempty"`
	Comment           string                     `xml:"comment,omitempty"`
	CreationTime      string                     `xml:"creation_time,omitempty"`     // ISO 8601 string
	ModificationTime  string                     `xml:"modification_time,omitempty"` // ISO 8601 string
	Writable          string                     `xml:"writable,omitempty"`          // "0"/"1"
	InUse             string                     `xml:"in_use,omitempty"`            // "0"/"1"
	Permissions       *TLSCertificatePermissions `xml:"permissions,omitempty"`
	Certificate       *TLSCertificateData        `xml:"certificate,omitempty"`
	SHA256Fingerprint string                     `xml:"sha256_fingerprint,omitempty"`
	MD5Fingerprint    string                     `xml:"md5_fingerprint,omitempty"`
	Trust             string                     `xml:"trust,omitempty"` // "0"/"1"
	Valid             string                     `xml:"valid,omitempty"` // "0"/"1"
	TimeStatus        string                     `xml:"time_status,omitempty"`
	ActivationTime    string                     `xml:"activation_time,omitempty"` // ISO 8601 string
	ExpirationTime    string                     `xml:"expiration_time,omitempty"` // ISO 8601 string
	SubjectDN         string                     `xml:"subject_dn,omitempty"`
	IssuerDN          string                     `xml:"issuer_dn,omitempty"`
	Serial            string                     `xml:"serial,omitempty"`
	LastSeen          string                     `xml:"last_seen,omitempty"` // ISO 8601 string
}

// TLSCertificateOwner represents the <owner> element for a TLS certificate.
type TLSCertificateOwner struct {
	Name string `xml:"name"`
}

// TLSCertificatePermissions represents the <permissions> element for a TLS certificate.
type TLSCertificatePermissions struct {
	Permission []TLSCertificatePermission `xml:"permission"`
}

// TLSCertificatePermission represents a <permission> element for a TLS certificate.
type TLSCertificatePermission struct {
	Name string `xml:"name"`
}

// TLSCertificateData represents the <certificate> element for a TLS certificate.
type TLSCertificateData struct {
	Format string `xml:"format,attr"`
	Data   string `xml:",chardata"`
}
