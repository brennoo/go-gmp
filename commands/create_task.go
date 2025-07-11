package commands

import "encoding/xml"

// CreateTask represents a create_task command request.
type CreateTask struct {
	XMLName         xml.Name         `xml:"create_task"`
	Name            string           `xml:"name,omitempty"`
	Comment         string           `xml:"comment,omitempty"`
	Copy            string           `xml:"copy,omitempty"`
	Alterable       bool             `xml:"alterable,omitempty"`
	UsageType       string           `xml:"usage_type,omitempty"`
	Config          *TaskConfig      `xml:"config,omitempty"`
	Target          *TaskTarget      `xml:"target,omitempty"`
	HostsOrdering   string           `xml:"hosts_ordering,omitempty"`
	Scanner         *TaskScanner     `xml:"scanner,omitempty"`
	Alert           *TaskAlert       `xml:"alert,omitempty"`
	Schedule        *TaskSchedule    `xml:"schedule,omitempty"`
	SchedulePeriods int              `xml:"schedule_periods,omitempty"`
	Observers       *TaskObservers   `xml:"observers,omitempty"`
	Preferences     *TaskPreferences `xml:"preferences,omitempty"`
}

// CreateTaskResponse represents a create_task command response.
type CreateTaskResponse struct {
	XMLName    xml.Name `xml:"create_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
