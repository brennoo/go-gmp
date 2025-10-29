package commands

import "encoding/xml"

// GetFeatures represents a get_features command request.
type GetFeatures struct {
	XMLName xml.Name `xml:"get_features"`
}

// GetFeaturesResponse represents a get_features command response.
type GetFeaturesResponse struct {
	XMLName    xml.Name  `xml:"get_features_response"`
	Status     string    `xml:"status,attr"`
	StatusText string    `xml:"status_text,attr"`
	Features   []Feature `xml:"feature"`
}

// GetStatus returns the status and status text from the response.
func (r *GetFeaturesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Feature represents a <feature> element in the get_features response.
type Feature struct {
	Enabled bool   `xml:"enabled,attr"`
	Name    string `xml:"name"`
}
