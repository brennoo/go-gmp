package gmp

import "encoding/xml"

// CreateAssetCommand represents a create_asset command request.
type CreateAssetCommand struct {
	XMLName xml.Name           `xml:"create_asset"`
	Asset   *CreateAssetAsset  `xml:"asset,omitempty"`
	Report  *CreateAssetReport `xml:"report,omitempty"`
}

// CreateAssetAsset represents the <asset> element for create_asset.
type CreateAssetAsset struct {
	Name    string `xml:"name"` // Must be an IPv4 or IPv6 address for hosts
	Comment string `xml:"comment,omitempty"`
	Type    string `xml:"type"` // Must be 'host'
}

// CreateAssetReport represents the <report> element for create_asset (import assets from report).
type CreateAssetReport struct {
	ID     string                   `xml:"id,attr"`
	Filter *CreateAssetReportFilter `xml:"filter,omitempty"`
}

// CreateAssetReportFilter represents the <filter> element for create_asset report import.
type CreateAssetReportFilter struct {
	Term string `xml:"term"`
}

// CreateAssetResponse represents a create_asset command response.
type CreateAssetResponse struct {
	XMLName    xml.Name `xml:"create_asset_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
