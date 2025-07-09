package commands

import (
	"encoding/xml"
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

// ScheduleOwner represents the owner of a schedule.
type ScheduleOwner struct {
	Name string `xml:"name"`
}

// SchedulePermission represents a permission on a schedule.
type SchedulePermission struct {
	Name string `xml:"name"`
}

// ScheduleUserTag represents a tag attached to a schedule.
type ScheduleUserTag struct {
	ID      string `xml:"id,attr"`
	Name    string `xml:"name"`
	Value   string `xml:"value"`
	Comment string `xml:"comment"`
}

// ScheduleUserTags represents tags attached to a schedule.
type ScheduleUserTags struct {
	Count int               `xml:"count"`
	Tags  []ScheduleUserTag `xml:"tag,omitempty"`
}

// ScheduleTask represents a task using a schedule.
type ScheduleTask struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name"`
	Permissions string `xml:"permissions,omitempty"`
}

// Schedule represents a schedule.
type Schedule struct {
	ID               string               `xml:"id,attr"`
	Owner            ScheduleOwner        `xml:"owner"`
	Name             string               `xml:"name"`
	Comment          string               `xml:"comment"`
	CreationTime     string               `xml:"creation_time"`
	ModificationTime string               `xml:"modification_time"`
	Writable         string               `xml:"writable"`
	InUse            string               `xml:"in_use"`
	Permissions      []SchedulePermission `xml:"permissions>permission,omitempty"`
	UserTags         *ScheduleUserTags    `xml:"user_tags,omitempty"`
	ICalendar        string               `xml:"icalendar"`
	Timezone         string               `xml:"timezone"`
	Tasks            []ScheduleTask       `xml:"tasks>task,omitempty"`
}

// GetSchedulesResponse represents a get_schedules command response.
type GetSchedulesResponse struct {
	XMLName       xml.Name   `xml:"get_schedules_response"`
	Status        string     `xml:"status,attr"`
	StatusText    string     `xml:"status_text,attr"`
	Schedules     []Schedule `xml:"schedule,omitempty"`
	Filters       string     `xml:"filters,omitempty"`
	Sort          string     `xml:"sort,omitempty"`
	ScheduleInfo  string     `xml:"schedules,omitempty"`
	ScheduleCount string     `xml:"schedule_count,omitempty"`
}
