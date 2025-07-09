package commands

import "encoding/xml"

// GetTLSCertificates represents a get_tls_certificates command request.
type GetTLSCertificates struct {
	XMLName                xml.Name `xml:"get_tls_certificates"`
	IncludeCertificateData string   `xml:"include_certificate_data,attr,omitempty"` // protocol: boolean as string
}

type TLSCertificate struct {
	ID                string           `xml:"id,attr"`
	Owner             *Owner           `xml:"owner,omitempty"`
	Name              string           `xml:"name,omitempty"`
	Comment           string           `xml:"comment,omitempty"`
	CreationTime      string           `xml:"creation_time,omitempty"`
	ModificationTime  string           `xml:"modification_time,omitempty"`
	Writable          string           `xml:"writable,omitempty"`
	InUse             string           `xml:"in_use,omitempty"`
	Permissions       *Permissions     `xml:"permissions,omitempty"`
	Certificate       *CertificateData `xml:"certificate,omitempty"`
	SHA256Fingerprint string           `xml:"sha256_fingerprint,omitempty"`
	MD5Fingerprint    string           `xml:"md5_fingerprint,omitempty"`
	Trust             string           `xml:"trust,omitempty"`
	Valid             string           `xml:"valid,omitempty"`
	TimeStatus        string           `xml:"time_status,omitempty"`
	ActivationTime    string           `xml:"activation_time,omitempty"`
	ExpirationTime    string           `xml:"expiration_time,omitempty"`
	SubjectDN         string           `xml:"subject_dn,omitempty"`
	IssuerDN          string           `xml:"issuer_dn,omitempty"`
	Serial            string           `xml:"serial,omitempty"`
	LastSeen          string           `xml:"last_seen,omitempty"`
}

type Owner struct {
	Name string `xml:"name"`
}

type Permissions struct {
	Permission []Permission `xml:"permission"`
}

type Permission struct {
	Name string `xml:"name"`
}

type CertificateData struct {
	Format string `xml:"format,attr"`
	Data   string `xml:",chardata"`
}

// GetTLSCertificatesResponse represents a get_tls_certificates command response.
type GetTLSCertificatesResponse struct {
	XMLName         xml.Name         `xml:"get_tls_certificates_response"`
	Status          string           `xml:"status,attr"`
	StatusText      string           `xml:"status_text,attr"`
	TLSCertificates []TLSCertificate `xml:"tls_certificate"`
}
