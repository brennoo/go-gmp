package commands

import "encoding/xml"

// ModifyLicenseCommand represents a modify_license command request.
type ModifyLicenseCommand struct {
	XMLName    xml.Name `xml:"modify_license"`
	AllowEmpty string   `xml:"allow_empty,attr,omitempty"`
	File       string   `xml:"file"`
}

// ModifyLicenseResponse represents a modify_license command response.
type ModifyLicenseResponse struct {
	XMLName    xml.Name `xml:"modify_license_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
