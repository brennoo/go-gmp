package commands

import (
	"encoding/xml"
	"time"
)

// GetTLSCertificates represents a get_tls_certificates command request.
type GetTLSCertificates struct {
	XMLName                xml.Name `xml:"get_tls_certificates"`
	TLSCertificateID       string   `xml:"tls_certificate_id,attr,omitempty"`
	Filter                 string   `xml:"filter,attr,omitempty"`
	IncludeCertificateData bool     `xml:"include_certificate_data,attr,omitempty"`
}

// GetTLSCertificatesResponse represents a get_tls_certificates command response.
type GetTLSCertificatesResponse struct {
	XMLName             xml.Name               `xml:"get_tls_certificates_response"`
	Status              string                 `xml:"status,attr"`
	StatusText          string                 `xml:"status_text,attr"`
	TLSCertificates     []TLSCertificate       `xml:"tls_certificate"`
	Filters             *TLSCertificateFilters `xml:"filters,omitempty"`
	Sort                *TLSCertificateSort    `xml:"sort,omitempty"`
	TLSCertificatesInfo *TLSCertificateTLS     `xml:"tls_certificates,omitempty"`
	TLSCertificateCount *TLSCertificateCount   `xml:"tls_certificate_count,omitempty"`
	Truncated           string                 `xml:"truncated,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetTLSCertificatesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// TLSCertificate represents a <tls_certificate> element in the get_tls_certificates response.
type TLSCertificate struct {
	ID                string                     `xml:"id,attr"`
	Owner             *TLSCertificateOwner       `xml:"owner,omitempty"`
	Name              string                     `xml:"name,omitempty"`
	Comment           string                     `xml:"comment,omitempty"`
	CreationTime      time.Time                  `xml:"creation_time,omitempty"`
	ModificationTime  time.Time                  `xml:"modification_time,omitempty"`
	Writable          bool                       `xml:"writable,omitempty"`
	InUse             bool                       `xml:"in_use,omitempty"`
	Permissions       *TLSCertificatePermissions `xml:"permissions,omitempty"`
	Certificate       *TLSCertificateData        `xml:"certificate,omitempty"`
	SHA256Fingerprint string                     `xml:"sha256_fingerprint,omitempty"`
	MD5Fingerprint    string                     `xml:"md5_fingerprint,omitempty"`
	Trust             bool                       `xml:"trust,omitempty"`
	Valid             bool                       `xml:"valid,omitempty"`
	TimeStatus        string                     `xml:"time_status,omitempty"`
	ActivationTime    time.Time                  `xml:"activation_time,omitempty"`
	ExpirationTime    time.Time                  `xml:"expiration_time,omitempty"`
	SubjectDN         string                     `xml:"subject_dn,omitempty"`
	IssuerDN          string                     `xml:"issuer_dn,omitempty"`
	Serial            string                     `xml:"serial,omitempty"`
	LastSeen          time.Time                  `xml:"last_seen,omitempty"`
	Sources           *TLSCertificateSources     `xml:"sources,omitempty"`
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

// TLSCertificateSources represents the <sources> element for a TLS certificate.
type TLSCertificateSources struct {
	Source []TLSCertificateSource `xml:"source,omitempty"`
}

// TLSCertificateSource represents a <source> element in sources.
type TLSCertificateSource struct {
	Type string `xml:"type,attr"`
	ID   string `xml:"id,attr"`
}

// TLSCertificateFilters represents the <filters> element in the response.
type TLSCertificateFilters struct {
	ID       string                         `xml:"id,attr,omitempty"`
	Term     string                         `xml:"term,omitempty"`
	Name     string                         `xml:"name,omitempty"`
	Comment  string                         `xml:"comment,omitempty"`
	Keywords *TLSCertificateFiltersKeywords `xml:"keywords,omitempty"`
}

// TLSCertificateFiltersKeywords represents the <keywords> element in filters.
type TLSCertificateFiltersKeywords struct {
	Keyword []TLSCertificateFiltersKeywordsKeyword `xml:"keyword,omitempty"`
}

// TLSCertificateFiltersKeywordsKeyword represents a <keyword> element in filters keywords.
type TLSCertificateFiltersKeywordsKeyword struct {
	Column   string `xml:"column,omitempty"`
	Relation string `xml:"relation,omitempty"`
	Value    string `xml:"value,omitempty"`
}

// TLSCertificateSort represents the <sort> element in the response.
type TLSCertificateSort struct {
	Value string                   `xml:",chardata"`
	Field *TLSCertificateSortField `xml:"field,omitempty"`
}

// TLSCertificateSortField represents the <field> element in sort.
type TLSCertificateSortField struct {
	Order string `xml:"order,omitempty"`
}

// TLSCertificateTLS represents the <tls_certificates> meta element in the response.
type TLSCertificateTLS struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// TLSCertificateCount represents the <tls_certificate_count> element in the response.
type TLSCertificateCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
