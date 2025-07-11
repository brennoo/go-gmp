package commands

import "encoding/xml"

// ModifyTask represents a modify_task command request.
type ModifyTask struct {
	XMLName         xml.Name         `xml:"modify_task"`
	TaskID          string           `xml:"task_id,attr"`
	Name            string           `xml:"name,omitempty"`
	Comment         string           `xml:"comment,omitempty"`
	Config          *TaskConfig      `xml:"config,omitempty"`
	Target          *TaskTarget      `xml:"target,omitempty"`
	Scanner         *TaskScanner     `xml:"scanner,omitempty"`
	Alert           []*TaskAlert     `xml:"alert,omitempty"`
	Schedule        *TaskSchedule    `xml:"schedule,omitempty"`
	SchedulePeriods int              `xml:"schedule_periods,omitempty"`
	Observers       string           `xml:"observers,omitempty"`
	Preferences     *TaskPreferences `xml:"preferences,omitempty"`
}

// ModifyTaskResponse represents a modify_task command response.
type ModifyTaskResponse struct {
	XMLName    xml.Name `xml:"modify_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
