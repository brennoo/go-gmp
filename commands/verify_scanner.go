package commands

import "encoding/xml"

// VerifyScannerCommand represents a verify_scanner command request.
type VerifyScannerCommand struct {
	XMLName   xml.Name `xml:"verify_scanner"`
	ScannerID string   `xml:"scanner_id,attr"`
}

// VerifyScannerResponse represents a verify_scanner command response.
type VerifyScannerResponse struct {
	XMLName    xml.Name `xml:"verify_scanner_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Version    string   `xml:"version,omitempty"`
}
