package commands

import "encoding/xml"

// GetLicense represents a get_license command request.
type GetLicense struct {
	XMLName xml.Name `xml:"get_license"`
}

// GetLicenseResponse represents a get_license command response.
type GetLicenseResponse struct {
	XMLName    xml.Name `xml:"get_license_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	License    *License `xml:"license,omitempty"`
}

// License represents the license element in the response.
type License struct {
	Status  string          `xml:"status"`
	Content *LicenseContent `xml:"content"`
}

// LicenseContent represents the content element in the license.
type LicenseContent struct {
	Meta       *LicenseMeta       `xml:"meta"`
	Appliance  *LicenseAppliance  `xml:"appliance"`
	Keys       *LicenseKeys       `xml:"keys"`
	Signatures *LicenseSignatures `xml:"signatures"`
}

// LicenseMeta represents the meta element in license content.
type LicenseMeta struct {
	Version      string `xml:"version"`
	ID           string `xml:"id"`
	Comment      string `xml:"comment"`
	Type         string `xml:"type"`
	CustomerName string `xml:"customer_name"`
	Created      string `xml:"created"`
	Begins       string `xml:"begins"`
	Expires      string `xml:"expires"`
}

// LicenseAppliance represents the appliance element in license content.
type LicenseAppliance struct {
	Model     string `xml:"model"`
	ModelType string `xml:"model_type"`
	Sensor    bool   `xml:"sensor"`
}

// LicenseKeys represents the keys element in license content.
type LicenseKeys struct {
	Keys []LicenseKey `xml:"key"`
}

// LicenseKey represents a key element in license keys.
type LicenseKey struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// LicenseSignatures represents the signatures element in license content.
type LicenseSignatures struct {
	Signatures []LicenseSignature `xml:"signature"`
}

// LicenseSignature represents a signature element in license signatures.
type LicenseSignature struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
