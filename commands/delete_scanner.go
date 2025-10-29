package commands

import "encoding/xml"

// DeleteScanner represents a delete_scanner command request.
type DeleteScanner struct {
	XMLName   xml.Name `xml:"delete_scanner"`
	ScannerID string   `xml:"scanner_id,attr"`
	Ultimate  bool     `xml:"ultimate,attr"`
}

// DeleteScannerResponse represents a delete_scanner command response.
type DeleteScannerResponse struct {
	XMLName    xml.Name `xml:"delete_scanner_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteScannerResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
