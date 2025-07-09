package commands

import "encoding/xml"

// ModifyTLSCertificate represents a modify_tls_certificate command request.
type ModifyTLSCertificate struct {
	XMLName          xml.Name `xml:"modify_tls_certificate"`
	TLSCertificateID string   `xml:"tls_certificate_id,attr"`
	Comment          string   `xml:"comment,omitempty"`
	Copy             string   `xml:"copy,omitempty"`
	Name             string   `xml:"name,omitempty"`
	Trust            string   `xml:"trust,omitempty"` // protocol: boolean as string
}

// ModifyTLSCertificateResponse represents a modify_tls_certificate command response.
type ModifyTLSCertificateResponse struct {
	XMLName    xml.Name `xml:"modify_tls_certificate_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
