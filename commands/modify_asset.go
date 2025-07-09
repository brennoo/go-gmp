package commands

import "encoding/xml"

// ModifyAssetCommand represents a modify_asset command request.
type ModifyAssetCommand struct {
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
