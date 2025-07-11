package commands

import "encoding/xml"

// CreateScanner is a request to create a scanner.
type CreateScanner struct {
	XMLName    xml.Name           `xml:"create_scanner"`
	Name       string             `xml:"name"`
	Comment    string             `xml:"comment,omitempty"`
	Copy       string             `xml:"copy,omitempty"`
	Host       string             `xml:"host"`
	Port       string             `xml:"port"`
	Type       string             `xml:"type"`
	CAPub      string             `xml:"ca_pub"`
	Credential *ScannerCredential `xml:"credential"`
	RelayHost  string             `xml:"relay_host,omitempty"`
	RelayPort  string             `xml:"relay_port,omitempty"`
}

// ScannerCredential is reused from get_scanners.go for both create and get scanner commands.
type CreateScannerResponse struct {
	XMLName    xml.Name `xml:"create_scanner_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
