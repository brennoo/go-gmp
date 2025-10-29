package commands

import (
	"encoding/xml"
	"time"
)

// GetAssets represents a get_assets command request.
type GetAssets struct {
	XMLName          xml.Name `xml:"get_assets"`
	AssetID          string   `xml:"asset_id,attr,omitempty"`
	Filter           string   `xml:"filter,attr,omitempty"`
	FiltID           string   `xml:"filt_id,attr,omitempty"`
	IgnorePagination string   `xml:"ignore_pagination,attr,omitempty"`
	Type             string   `xml:"type,attr,omitempty"`
	Details          string   `xml:"details,attr,omitempty"`
}

// GetAssetsResponse represents a get_assets command response.
type GetAssetsResponse struct {
	XMLName    xml.Name      `xml:"get_assets_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Assets     []Asset       `xml:"asset,omitempty"`
	Filters    *AssetFilters `xml:"filters,omitempty"`
	Sort       *AssetSort    `xml:"sort,omitempty"`
	AssetsInfo *AssetAssets  `xml:"assets,omitempty"`
	AssetCount *AssetCount   `xml:"asset_count,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetAssetsResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Asset represents a single asset in the response.
type Asset struct {
	ID               string            `xml:"id,attr"`
	Owner            AssetOwner        `xml:"owner"`
	Name             string            `xml:"name"`
	Comment          string            `xml:"comment"`
	CreationTime     time.Time         `xml:"creation_time"`
	ModificationTime time.Time         `xml:"modification_time"`
	Writable         bool              `xml:"writable,omitempty"`
	InUse            bool              `xml:"in_use,omitempty"`
	Permissions      []AssetPermission `xml:"permissions>permission,omitempty"`
	UserTags         *AssetUserTags    `xml:"user_tags,omitempty"`
	Identifiers      *AssetIdentifiers `xml:"identifiers,omitempty"`
	Type             string            `xml:"type"`
	Hosts            string            `xml:"hosts,omitempty"`
	MaxHosts         int               `xml:"max_hosts,omitempty"`
	SSHCredential    *AssetCredential  `xml:"ssh_credential,omitempty"`
	SMBCredential    *AssetCredential  `xml:"smb_credential,omitempty"`
	ESXICredential   *AssetCredential  `xml:"esxi_credential,omitempty"`
	Host             *AssetHost        `xml:"host,omitempty"`
	OS               *AssetOS          `xml:"os,omitempty"`
	Tasks            *AssetTasks       `xml:"tasks,omitempty"`
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

// AssetIdentifier represents a single identifier in an asset.
type AssetIdentifier struct {
	Name             string                `xml:"name"`
	Value            string                `xml:"value"`
	CreationTime     string                `xml:"creation_time"`
	ModificationTime string                `xml:"modification_time"`
	Source           AssetIdentifierSource `xml:"source"`
	OS               *AssetIdentifierOS    `xml:"os,omitempty"`
}

// AssetIdentifierSource represents the source of an identifier.
type AssetIdentifierSource struct {
	ID      string `xml:"id,attr"`
	Type    string `xml:"type"`
	Data    string `xml:"data"`
	Deleted string `xml:"deleted"`
	Name    string `xml:"name"`
}

// AssetIdentifierOS represents an OS in an identifier.
type AssetIdentifierOS struct {
	ID    string `xml:"id,attr"`
	Title string `xml:"title"`
}

// AssetIdentifiers represents identifiers in an asset.
type AssetIdentifiers struct {
	Identifier []AssetIdentifier `xml:"identifier"`
}

// AssetHost represents a host asset.
type AssetHost struct {
	Severity AssetHostSeverity `xml:"severity"`
	Details  []AssetHostDetail `xml:"detail,omitempty"`
	Routes   *AssetHostRoutes  `xml:"routes,omitempty"`
}

// AssetHostSeverity represents the severity of a host.
type AssetHostSeverity struct {
	Value string `xml:"value"`
}

// AssetHostDetail represents a detail of a host.
type AssetHostDetail struct {
	Name   string            `xml:"name"`
	Value  string            `xml:"value"`
	Source AssetDetailSource `xml:"source"`
}

// AssetDetailSource represents the source of a detail.
type AssetDetailSource struct {
	ID   string `xml:"id,attr"`
	Type string `xml:"type"`
}

// AssetHostRoute represents a route to a host.
type AssetHostRoute struct {
	Hosts []AssetRouteHost `xml:"host,omitempty"`
}

// AssetRouteHost represents a host on a route.
type AssetRouteHost struct {
	ID         string `xml:"id,attr"`
	Distance   string `xml:"distance,attr"`
	SameSource string `xml:"same_source,attr"`
	IP         string `xml:"ip"`
}

// AssetHostRoutes represents routes to a host.
type AssetHostRoutes struct {
	Routes []AssetHostRoute `xml:"route,omitempty"`
}

// AssetOS represents an OS asset.
type AssetOS struct {
	Title           string        `xml:"title"`
	Installs        int           `xml:"installs"`
	AllInstalls     int           `xml:"all_installs"`
	LatestSeverity  AssetSeverity `xml:"latest_severity"`
	HighestSeverity AssetSeverity `xml:"highest_severity"`
	AverageSeverity AssetSeverity `xml:"average_severity"`
	Hosts           []AssetOSHost `xml:"hosts>asset,omitempty"`
}

// AssetSeverity represents a severity value.
type AssetSeverity struct {
	Value string `xml:"value"`
}

// AssetOSHost represents a host in an OS asset.
type AssetOSHost struct {
	ID       string        `xml:"id,attr"`
	Name     string        `xml:"name"`
	Severity AssetSeverity `xml:"severity"`
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

// AssetFilters represents filter information in the response.
type AssetFilters struct {
	ID       string        `xml:"id,attr"`
	Term     string        `xml:"term"`
	Name     string        `xml:"name,omitempty"`
	Keywords AssetKeywords `xml:"keywords"`
}

// AssetKeyword represents a keyword in filters.
type AssetKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// AssetKeywords represents keywords in filters.
type AssetKeywords struct {
	Keyword []AssetKeyword `xml:"keyword"`
}

// AssetSort represents sorting information in the response.
type AssetSort struct {
	Text  string         `xml:",chardata"`
	Field AssetSortField `xml:"field"`
}

// AssetSortField represents the field element in sort.
type AssetSortField struct {
	Order string `xml:"order"`
}

// AssetAssets represents the assets container in the response.
type AssetAssets struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// AssetCount represents count information in the response.
type AssetCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}

// AssetCredential represents a credential associated with an asset.
type AssetCredential struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}
