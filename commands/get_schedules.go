package commands

import (
	"encoding/xml"
	"time"
)

// GetSchedules represents a get_schedules command request.
type GetSchedules struct {
	XMLName    xml.Name `xml:"get_schedules"`
	ScheduleID string   `xml:"schedule_id,attr,omitempty"`
	Filter     string   `xml:"filter,attr,omitempty"`
	FiltID     string   `xml:"filt_id,attr,omitempty"`
	Trash      string   `xml:"trash,attr,omitempty"`
	Tasks      string   `xml:"tasks,attr,omitempty"`
}

// GetSchedulesResponse represents a get_schedules command response.
type GetSchedulesResponse struct {
	XMLName       xml.Name         `xml:"get_schedules_response"`
	Status        string           `xml:"status,attr"`
	StatusText    string           `xml:"status_text,attr"`
	Schedules     []Schedule       `xml:"schedule,omitempty"`
	Filters       *ScheduleFilters `xml:"filters,omitempty"`
	Sort          *ScheduleSort    `xml:"sort,omitempty"`
	ScheduleInfo  *ScheduleInfo    `xml:"schedules,omitempty"`
	ScheduleCount *ScheduleCount   `xml:"schedule_count,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetSchedulesResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Schedule represents a <schedule> element in the get_schedules response.
type Schedule struct {
	ID               string               `xml:"id,attr"`
	Owner            ScheduleOwner        `xml:"owner"`
	Name             string               `xml:"name"`
	Comment          string               `xml:"comment"`
	CreationTime     time.Time            `xml:"creation_time"`
	ModificationTime time.Time            `xml:"modification_time"`
	Writable         bool                 `xml:"writable"`
	InUse            bool                 `xml:"in_use"`
	Permissions      []SchedulePermission `xml:"permissions>permission,omitempty"`
	UserTags         *ScheduleUserTags    `xml:"user_tags,omitempty"`
	ICalendar        string               `xml:"icalendar"`
	Timezone         string               `xml:"timezone"`
	Tasks            []ScheduleTask       `xml:"tasks>task,omitempty"`
}

// ScheduleOwner represents the owner of a schedule.
type ScheduleOwner struct {
	Name string `xml:"name"`
}

// SchedulePermission represents a permission on a schedule.
type SchedulePermission struct {
	Name string `xml:"name"`
}

// ScheduleUserTags represents tags attached to a schedule.
type ScheduleUserTags struct {
	Count int               `xml:"count"`
	Tags  []ScheduleUserTag `xml:"tag,omitempty"`
}

// ScheduleUserTag represents a tag attached to a schedule.
type ScheduleUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ScheduleTask represents a task using a schedule.
type ScheduleTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
}

// ScheduleFilters represents the filters element in the response.
type ScheduleFilters struct {
	ID       string            `xml:"id,attr,omitempty"`
	Term     string            `xml:"term,omitempty"`
	Name     string            `xml:"name,omitempty"`
	Keywords *ScheduleKeywords `xml:"keywords,omitempty"`
}

// ScheduleKeywords represents the keywords element in filters.
type ScheduleKeywords struct {
	Keyword []ScheduleKeyword `xml:"keyword,omitempty"`
}

// ScheduleKeyword represents a keyword element in filters.
type ScheduleKeyword struct {
	Column   string `xml:"column,omitempty"`
	Relation string `xml:"relation,omitempty"`
	Value    string `xml:"value,omitempty"`
}

// ScheduleSort represents the sort element in the response.
type ScheduleSort struct {
	Value string         `xml:",chardata"`
	Field *ScheduleField `xml:"field,omitempty"`
}

// ScheduleField represents a field element in sort.
type ScheduleField struct {
	Value string `xml:",chardata"`
	Order string `xml:"order,omitempty"`
}

// ScheduleInfo represents the schedules element in the response.
type ScheduleInfo struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// ScheduleCount represents the schedule_count element in the response.
type ScheduleCount struct {
	Filtered int `xml:"filtered,omitempty"`
	Page     int `xml:"page,omitempty"`
}
