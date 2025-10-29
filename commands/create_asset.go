package commands

import "encoding/xml"

// CreateAsset represents a create_asset command request.
type CreateAsset struct {
	XMLName xml.Name     `xml:"create_asset"`
	Asset   *AssetInput  `xml:"asset,omitempty"`
	Report  *AssetReport `xml:"report,omitempty"`
}

// CreateAssetResponse represents a create_asset command response.
type CreateAssetResponse struct {
	XMLName    xml.Name `xml:"create_asset_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateAssetResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// AssetInput represents the <asset> element for create_asset.
type AssetInput struct {
	Name    string `xml:"name"`
	Comment string `xml:"comment,omitempty"`
	Type    string `xml:"type"`
}

// AssetReport represents the <report> element for create_asset.
type AssetReport struct {
	ID     string       `xml:"id,attr"`
	Filter *AssetFilter `xml:"filter,omitempty"`
}

// CreateAssetFilter represents the <filter> element for create_asset report.
type AssetFilter struct {
	Term string `xml:"term"`
}
