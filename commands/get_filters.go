package commands

import "encoding/xml"

// GetFilters represents a get_filters command request.
type GetFilters struct {
	XMLName  xml.Name `xml:"get_filters"`
	FilterID string   `xml:"filter_id,attr,omitempty"`
	Filter   string   `xml:"filter,attr,omitempty"`
	FiltID   string   `xml:"filt_id,attr,omitempty"`
	Trash    bool     `xml:"trash,attr,omitempty"`
	Alerts   bool     `xml:"alerts,attr,omitempty"`
}

// GetFiltersResponse represents a get_filters command response.
type GetFiltersResponse struct {
	XMLName     xml.Name    `xml:"get_filters_response"`
	Status      string      `xml:"status,attr"`
	StatusText  string      `xml:"status_text,attr"`
	Filters     []Filter    `xml:"filter"`
	FiltersInfo FilterInfo  `xml:"filters"`
	Sort        FilterSort  `xml:"sort"`
	Filters2    FilterInfo2 `xml:"filters2"`
	FilterCount FilterCount `xml:"filter_count"`
}

// Filter represents a <filter> element in the get_filters response.
type Filter struct {
	ID               string            `xml:"id,attr"`
	Owner            FilterOwner       `xml:"owner"`
	Name             string            `xml:"name"`
	Comment          string            `xml:"comment"`
	Term             string            `xml:"term"`
	Type             string            `xml:"type"`
	CreationTime     string            `xml:"creation_time"`
	ModificationTime string            `xml:"modification_time"`
	InUse            bool              `xml:"in_use"`
	Writable         bool              `xml:"writable"`
	Permissions      FilterPermissions `xml:"permissions"`
	UserTags         *FilterUserTags   `xml:"user_tags,omitempty"`
	Alerts           *FilterAlerts     `xml:"alerts,omitempty"`
}

// FilterOwner represents the owner element in a filter.
type FilterOwner struct {
	Name string `xml:"name"`
}

// FilterPermissions represents the permissions element in a filter.
type FilterPermissions struct {
	Permission []FilterPermission `xml:"permission"`
}

// FilterPermission represents a permission element in filter permissions.
type FilterPermission struct {
	Name string `xml:"name"`
}

// FilterUserTags represents the user tags element in a filter.
type FilterUserTags struct {
	Count int             `xml:"count"`
	Tags  []FilterUserTag `xml:"tag,omitempty"`
}

// FilterUserTag represents a tag element in filter user tags.
type FilterUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// FilterAlerts represents the alerts element in a filter.
type FilterAlerts struct {
	Alert []FilterAlert `xml:"alert"`
}

// FilterAlert represents an alert element in filter alerts.
type FilterAlert struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

// FilterInfo represents the filters element in the response.
type FilterInfo struct {
	ID       string         `xml:"id,attr"`
	Term     string         `xml:"term"`
	Name     string         `xml:"name"`
	Keywords FilterKeywords `xml:"keywords"`
}

// FilterKeywords represents the keywords element in filters.
type FilterKeywords struct {
	Keyword []FilterKeyword `xml:"keyword"`
}

// FilterKeyword represents a keyword element in filters keywords.
type FilterKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// FilterSort represents the sort element in the response.
type FilterSort struct {
	Value string          `xml:",chardata"`
	Field FilterSortField `xml:"field"`
}

// FilterSortField represents the field element in sort.
type FilterSortField struct {
	Order string `xml:"order"`
}

// FilterInfo2 represents the filters2 element in the response.
type FilterInfo2 struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// FilterCount represents the filter count element in the response.
type FilterCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
