package commands

import "encoding/xml"

// GetFeaturesCommand represents a get_features command request.
type GetFeaturesCommand struct {
	XMLName xml.Name `xml:"get_features"`
}

// GetFeaturesResponse represents a get_features command response.
type GetFeaturesResponse struct {
	XMLName    xml.Name      `xml:"get_features_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Features   []FeatureInfo `xml:"feature"`
}

type FeatureInfo struct {
	Enabled string `xml:"enabled,attr"`
	Name    string `xml:"name"`
}
