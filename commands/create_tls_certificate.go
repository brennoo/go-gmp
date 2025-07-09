package commands

import "encoding/xml"

// CreateTLSCertificateCommand represents a create_tls_certificate command request.
type CreateTLSCertificateCommand struct {
	XMLName     xml.Name `xml:"create_tls_certificate"`
	Comment     string   `xml:"comment,omitempty"`
	Copy        string   `xml:"copy,omitempty"`
	Name        string   `xml:"name,omitempty"`
	Trust       string   `xml:"trust,omitempty"` // protocol: boolean as string
	Certificate string   `xml:"certificate"`
}

// CreateTLSCertificateResponse represents a create_tls_certificate command response.
type CreateTLSCertificateResponse struct {
	XMLName    xml.Name `xml:"create_tls_certificate_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
