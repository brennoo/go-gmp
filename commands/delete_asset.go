package commands

import "encoding/xml"

// DeleteAsset represents a delete_asset command request.
type DeleteAsset struct {
	XMLName  xml.Name `xml:"delete_asset"`
	AssetID  string   `xml:"asset_id,attr,omitempty"`
	ReportID string   `xml:"report_id,attr,omitempty"`
}

// DeleteAssetResponse represents a delete_asset command response.
type DeleteAssetResponse struct {
	XMLName    xml.Name `xml:"delete_asset_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteAssetResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
