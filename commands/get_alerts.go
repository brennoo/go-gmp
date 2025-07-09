package commands

import (
	"encoding/xml"
)

// GetAlertsCommand represents a get_alerts command request.
type GetAlertsCommand struct {
	XMLName xml.Name `xml:"get_alerts"`
	AlertID string   `xml:"alert_id,attr,omitempty"`
	Filter  string   `xml:"filter,attr,omitempty"`
	FiltID  string   `xml:"filt_id,attr,omitempty"`
	Trash   *bool    `xml:"trash,attr,omitempty"`
	Tasks   *bool    `xml:"tasks,attr,omitempty"`
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

// FilterKeyword represents a keyword in a filter.
type FilterKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

// Filters represents filter information in the response.
type Filters struct {
	ID       string          `xml:"id,attr"`
	Term     string          `xml:"term"`
	Name     string          `xml:"name,omitempty"`
	Keywords []FilterKeyword `xml:"keywords>keyword,omitempty"`
}

// Sort represents sorting information in the response.
type Sort struct {
	Text  string `xml:",chardata"`
	Field struct {
		Order string `xml:"order"`
	} `xml:"field"`
}

// Alerts represents the alerts container in the response.
type Alerts struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

// AlertCount represents count information in the response.
type AlertCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}

// GetAlertsResponse represents a get_alerts command response.
type GetAlertsResponse struct {
	XMLName    xml.Name    `xml:"get_alerts_response"`
	Status     string      `xml:"status,attr"`
	StatusText string      `xml:"status_text,attr"`
	Alerts     []Alert     `xml:"alert,omitempty"`
	Filters    *Filters    `xml:"filters,omitempty"`
	Sort       *Sort       `xml:"sort,omitempty"`
	AlertsInfo *Alerts     `xml:"alerts,omitempty"`
	AlertCount *AlertCount `xml:"alert_count,omitempty"`
}
