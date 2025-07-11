package commands

import (
	"encoding/xml"
)

// GetAlerts represents a get_alerts command request.
type GetAlerts struct {
	XMLName xml.Name `xml:"get_alerts"`
	AlertID string   `xml:"alert_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
	FiltID  string   `xml:"filt_id,attr,omitempty"`
	Trash   *bool    `xml:"trash,attr,omitempty"`
	Tasks   *bool    `xml:"tasks,attr,omitempty"`
}

// GetAlertsResponse represents a get_alerts command response.
type GetAlertsResponse struct {
	XMLName    xml.Name      `xml:"get_alerts_response"`
	Status     string        `xml:"status,attr"`
	StatusText string        `xml:"status_text,attr"`
	Alerts     []Alert       `xml:"alert,omitempty"`
	Filters    *AlertFilters `xml:"filters,omitempty"`
	Sort       *AlertSort    `xml:"sort,omitempty"`
	AlertsInfo *AlertInfo    `xml:"alerts,omitempty"`
	AlertCount *AlertCount   `xml:"alert_count,omitempty"`
}

// Alert represents a single alert in the response.
type Alert struct {
	ID               string            `xml:"id,attr"`
	Owner            AlertOwner        `xml:"owner"`
	Name             string            `xml:"name"`
	Comment          string            `xml:"comment"`
	CreationTime     string            `xml:"creation_time"`
	ModificationTime string            `xml:"modification_time"`
	InUse            *bool             `xml:"in_use"`
	Writable         *bool             `xml:"writable"`
	Permissions      []AlertPermission `xml:"permissions>permission,omitempty"`
	UserTags         *AlertUserTags    `xml:"user_tags,omitempty"`
	Condition        AlertCondition    `xml:"condition"`
	Event            AlertEvent        `xml:"event"`
	Method           AlertMethod       `xml:"method"`
	Filter           AlertFilterInfo   `xml:"filter"`
	Tasks            []AlertTask       `xml:"tasks>task,omitempty"`
	Active           *bool             `xml:"active"`
}

// AlertOwner represents the owner of an alert.
type AlertOwner struct {
	Name string `xml:"name"`
}

// AlertPermission represents a permission on an alert.
type AlertPermission struct {
	Name string `xml:"name"`
}

// AlertUserTag represents a tag attached to an alert.
type AlertUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// AlertUserTags represents user tags attached to an alert.
type AlertUserTags struct {
	Count int            `xml:"count"`
	Tags  []AlertUserTag `xml:"tag,omitempty"`
}

// AlertFilterInfo represents a filter applied to an alert.
type AlertFilterInfo struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
	Trash       *bool  `xml:"trash,omitempty"`
}

// AlertTask represents a task using an alert.
type AlertTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
}

// AlertFilterKeyword represents a keyword in a filter.
type AlertFilterKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// AlertFilters represents filter information in the response.
type AlertFilters struct {
	ID       string               `xml:"id,attr"`
	Term     string               `xml:"term"`
	Name     string               `xml:"name,omitempty"`
	Keywords []AlertFilterKeyword `xml:"keywords>keyword,omitempty"`
}

// AlertSort represents sorting information in the response.
type AlertSort struct {
	Text  string         `xml:",chardata"`
	Field AlertSortField `xml:"field"`
}

// AlertSortField represents the field element in sort.
type AlertSortField struct {
	Order string `xml:"order"`
}

// AlertInfo represents the alerts container in the response.
type AlertInfo struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// AlertCount represents count information in the response.
type AlertCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}
