package commands

import "encoding/xml"

// GetVersion represents a get_version command request.
type GetVersion struct {
	XMLName xml.Name `xml:"get_version"`
}

// GetVersionResponse represents a get_version command response.
type GetVersionResponse struct {
	XMLName    xml.Name `xml:"get_version_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	Version    string   `xml:"version"`
}
