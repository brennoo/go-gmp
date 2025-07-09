package commands

import "encoding/xml"

// ModifyScannerCommand is a request to modify a scanner.
type ModifyScannerCommand struct {
	XMLName    xml.Name                 `xml:"modify_scanner"`
	ScannerID  string                   `xml:"scanner_id,attr"`
	Comment    string                   `xml:"comment,omitempty"`
	Name       string                   `xml:"name,omitempty"`
	Host       string                   `xml:"host,omitempty"`
	Port       string                   `xml:"port,omitempty"`
	Type       string                   `xml:"type,omitempty"`
	CAPub      string                   `xml:"ca_pub,omitempty"`
	Credential *CreateScannerCredential `xml:"credential,omitempty"`
	RelayHost  string                   `xml:"relay_host,omitempty"`
	RelayPort  string                   `xml:"relay_port,omitempty"`
}

// ModifyScannerResponse is the response to a modify_scanner command.
type ModifyScannerResponse struct {
	XMLName    xml.Name `xml:"modify_scanner_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
