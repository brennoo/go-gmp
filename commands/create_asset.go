package commands

import "encoding/xml"

// CreateAsset represents a create_asset command request.
type CreateAsset struct {
	XMLName xml.Name `xml:"create_asset"`
	Asset   *asset   `xml:"asset,omitempty"`
	Report  *report  `xml:"report,omitempty"`
}

// CreateAssetResponse represents a create_asset command response.
type CreateAssetResponse struct {
	XMLName    xml.Name `xml:"create_asset_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// NewCreateAsset creates a new CreateAsset command.
func NewCreateAsset(asset *asset, report *report) *CreateAsset {
	return &CreateAsset{
		Asset:  asset,
		Report: report,
	}
}

// asset represents the <asset> element for create_asset.
type asset struct {
	Name    string `xml:"name"`
	Comment string `xml:"comment,omitempty"`
	Type    string `xml:"type"`
}

// NewAsset creates a new asset for create_asset.
func NewAsset(name, comment, typ string) *asset {
	return &asset{
		Name:    name,
		Comment: comment,
		Type:    typ,
	}
}

// report represents the <report> element for create_asset.
type report struct {
	ID     string  `xml:"id,attr"`
	Filter *Filter `xml:"filter,omitempty"`
}

// NewAssetReport creates a new report for create_asset.
func NewAssetReport(id string, filter *Filter) *report {
	return &report{
		ID:     id,
		Filter: filter,
	}
}
