package commands

import "encoding/xml"

// ModifyAsset represents a modify_asset command request.
type ModifyAsset struct {
	XMLName xml.Name `xml:"modify_asset"`
	AssetID string   `xml:"asset_id,attr"`
	Comment string   `xml:"comment,omitempty"`
}

// ModifyAssetResponse represents a modify_asset command response.
type ModifyAssetResponse struct {
	XMLName    xml.Name `xml:"modify_asset_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyAssetResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
