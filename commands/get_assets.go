package commands

import (
	"encoding/xml"
)

// GetAssets represents a get_assets command request.
type GetAssets struct {
	XMLName xml.Name `xml:"get_assets"`
	AssetID string   `xml:"asset_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
}

// AssetOwner represents the owner of an asset.
type AssetOwner struct {
	Name string `xml:"name"`
}

// AssetPermission represents a permission on an asset.
type AssetPermission struct {
	Name string `xml:"name"`
}

// AssetUserTag represents a tag attached to an asset.
type AssetUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// AssetUserTags represents user tags attached to an asset.
type AssetUserTags struct {
	Count int            `xml:"count"`
	Tags  []AssetUserTag `xml:"tag,omitempty"`
}

// AssetTask represents a task using an asset.
type AssetTask struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// AssetTasks represents tasks using an asset.
type AssetTasks struct {
	Tasks []AssetTask `xml:"task,omitempty"`
}

// Asset represents a single asset in the response.
type Asset struct {
	ID               string            `xml:"id,attr"`
	Owner            AssetOwner        `xml:"owner"`
	Name             string            `xml:"name"`
	Comment          string            `xml:"comment"`
	CreationTime     string            `xml:"creation_time"`
	ModificationTime string            `xml:"modification_time"`
	Writable         *bool             `xml:"writable,omitempty"`
	InUse            *bool             `xml:"in_use,omitempty"`
	Permissions      []AssetPermission `xml:"permissions>permission,omitempty"`
	UserTags         *AssetUserTags    `xml:"user_tags,omitempty"`
	Type             string            `xml:"type"`
	Tasks            *AssetTasks       `xml:"tasks,omitempty"`
}

// AssetFilters represents filter information in the response.
type AssetFilters struct {
	ID   string `xml:"id,attr"`
	Term string `xml:"term"`
}

// AssetSort represents sorting information in the response.
type AssetSort struct {
	Text  string `xml:",chardata"`
	Field struct {
		Order string `xml:"order"`
	} `xml:"field"`
}

// AssetsInfo represents the assets container in the response.
type AssetsInfo struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// AssetCount represents count information in the response.
type AssetCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}

// GetAssetsResponse represents a get_assets command response.
type GetAssetsResponse struct {
	XMLName    xml.Name      `xml:"get_assets_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Assets     []Asset       `xml:"asset,omitempty"`
	Filters    *AssetFilters `xml:"filters,omitempty"`
	Sort       *AssetSort    `xml:"sort,omitempty"`
	AssetsInfo *AssetsInfo   `xml:"assets,omitempty"`
	AssetCount *AssetCount   `xml:"asset_count,omitempty"`
}
