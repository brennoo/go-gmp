package gmp

import "encoding/xml"

// GetLicenseCommand represents a get_license command request.
type GetLicenseCommand struct {
	XMLName xml.Name `xml:"get_license"`
}

// GetLicenseResponse represents a get_license command response.
type GetLicenseResponse struct {
	XMLName    xml.Name `xml:"get_license_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	License    *License `xml:"license,omitempty"`
}

type License struct {
	Status  string          `xml:"status,omitempty"`
	Content *LicenseContent `xml:"content,omitempty"`
}

type LicenseContent struct {
	Meta       *LicenseMeta       `xml:"meta,omitempty"`
	Appliance  *LicenseAppliance  `xml:"appliance,omitempty"`
	Keys       *LicenseKeys       `xml:"keys,omitempty"`
	Signatures *LicenseSignatures `xml:"signatures,omitempty"`
}

type LicenseMeta struct {
	Version      string `xml:"version,omitempty"`
	ID           string `xml:"id,omitempty"`
	Comment      string `xml:"comment,omitempty"`
	Type         string `xml:"type,omitempty"`
	CustomerName string `xml:"customer_name,omitempty"`
	Created      string `xml:"created,omitempty"`
	Begins       string `xml:"begins,omitempty"`
	Expires      string `xml:"expires,omitempty"`
}

type LicenseAppliance struct {
	Model     string `xml:"model,omitempty"`
	ModelType string `xml:"model_type,omitempty"`
	Sensor    string `xml:"sensor,omitempty"`
}

type LicenseKeys struct {
	Keys []LicenseKey `xml:"key,omitempty"`
}

type LicenseKey struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:",chardata"`
}

type LicenseSignatures struct {
	Signatures []LicenseSignature `xml:"signature,omitempty"`
}

type LicenseSignature struct {
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:",chardata"`
}
