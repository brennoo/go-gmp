package commands

import (
	"encoding/xml"
)

// GetPermissions represents a get_permissions command request.
type GetPermissions struct {
	XMLName      xml.Name `xml:"get_permissions"`
	PermissionID string   `xml:"permission_id,attr,omitempty"`
	Filter       string   `xml:"filter,attr,omitempty"`
	FiltID       string   `xml:"filt_id,attr,omitempty"`
	Trash        string   `xml:"trash,attr,omitempty"`
}

// GetPermissionsResponse represents a get_permissions command response.
type GetPermissionsResponse struct {
	XMLName     xml.Name              `xml:"get_permissions_response"`
	Status      string                `xml:"status,attr"`
	StatusText  string                `xml:"status_text,attr"`
	Permissions []Permission          `xml:"permission,omitempty"`
	Filters     *PermissionFilters    `xml:"filters,omitempty"`
	Sort        *PermissionSort       `xml:"sort,omitempty"`
	Pagination  *PermissionPagination `xml:"permissions,omitempty"`
	Count       *PermissionCount      `xml:"permission_count,omitempty"`
}

// Permission represents a permission in the response.
type Permission struct {
	ID               string              `xml:"id,attr"`
	Owner            *PermissionOwner    `xml:"owner,omitempty"`
	Name             string              `xml:"name,omitempty"`
	Comment          string              `xml:"comment,omitempty"`
	CreationTime     string              `xml:"creation_time,omitempty"`
	ModificationTime string              `xml:"modification_time,omitempty"`
	Writable         string              `xml:"writable,omitempty"`
	InUse            string              `xml:"in_use,omitempty"`
	Permissions      *PermissionList     `xml:"permissions,omitempty"`
	UserTags         *PermissionUserTags `xml:"user_tags,omitempty"`
	Resource         *PermissionResource `xml:"resource,omitempty"`
	Subject          *PermissionSubject  `xml:"subject,omitempty"`
}

// PermissionOwner represents owner information in a permission.
type PermissionOwner struct {
	Name string `xml:"name,omitempty"`
}

// PermissionList represents a list of permissions.
type PermissionList struct {
	Permission []PermissionItem `xml:"permission,omitempty"`
}

// PermissionItem represents a single permission item.
type PermissionItem struct {
	Name string `xml:"name,omitempty"`
}

// PermissionUserTags represents user tags information in a permission.
type PermissionUserTags struct {
	Count int                 `xml:"count,omitempty"`
	Tag   []PermissionUserTag `xml:"tag,omitempty"`
}

// PermissionUserTag represents a user tag in a permission.
type PermissionUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name,omitempty"`
	Value   string `xml:"value,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// PermissionResource represents resource information in a permission.
type PermissionResource struct {
	ID          string          `xml:"id,attr"`
	Name        string          `xml:"name,omitempty"`
	Permissions *PermissionList `xml:"permissions,omitempty"`
	Type        string          `xml:"type,omitempty"`
	Trash       string          `xml:"trash,omitempty"`
	Deleted     string          `xml:"deleted,omitempty"`
}

// PermissionSubject represents subject information in a permission.
type PermissionSubject struct {
	ID          string          `xml:"id,attr"`
	Name        string          `xml:"name,omitempty"`
	Permissions *PermissionList `xml:"permissions,omitempty"`
	Type        string          `xml:"type,omitempty"`
	Trash       string          `xml:"trash,omitempty"`
	Deleted     string          `xml:"deleted,omitempty"`
}

// PermissionFilters represents filter information in the response.
type PermissionFilters struct {
	ID       string                     `xml:"id,attr"`
	Term     string                     `xml:"term"`
	Name     string                     `xml:"name,omitempty"`
	Keywords []PermissionFiltersKeyword `xml:"keywords>keyword,omitempty"`
}

// PermissionFiltersKeyword represents a keyword in permission filters.
type PermissionFiltersKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// PermissionSort represents sort information in the response.
type PermissionSort struct {
	Field PermissionSortField `xml:"field"`
}

// PermissionSortField represents a sort field in permission sorting.
type PermissionSortField struct {
	Order string `xml:"order"`
}

// PermissionPagination represents pagination information in the response.
type PermissionPagination struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// PermissionCount represents count information in the response.
type PermissionCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
