package commands

import "encoding/xml"

// ModifyLicense represents a modify_license command request.
type ModifyLicense struct {
	XMLName    xml.Name `xml:"modify_license"`
	AllowEmpty bool     `xml:"allow_empty,attr,omitempty"`
	File       string   `xml:"file"`
}

// ModifyLicenseResponse represents a modify_license command response.
type ModifyLicenseResponse struct {
	XMLName    xml.Name `xml:"modify_license_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyLicenseResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
