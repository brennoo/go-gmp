package gmp

import "encoding/xml"

// ModifyTaskCommand represents a modify_task command request.
type ModifyTaskCommand struct {
	XMLName         xml.Name               `xml:"modify_task"`
	TaskID          string                 `xml:"task_id,attr"`
	Name            string                 `xml:"name,omitempty"`
	Comment         string                 `xml:"comment,omitempty"`
	Config          *ModifyTaskConfig      `xml:"config,omitempty"`
	Target          *ModifyTaskTarget      `xml:"target,omitempty"`
	Scanner         *ModifyTaskScanner     `xml:"scanner,omitempty"`
	Alert           []*ModifyTaskAlert     `xml:"alert,omitempty"`
	Schedule        *ModifyTaskSchedule    `xml:"schedule,omitempty"`
	SchedulePeriods int                    `xml:"schedule_periods,omitempty"`
	Observers       string                 `xml:"observers,omitempty"`
	Preferences     *ModifyTaskPreferences `xml:"preferences,omitempty"`
	File            *ModifyTaskFile        `xml:"file,omitempty"`
}

type ModifyTaskConfig struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTaskTarget struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTaskScanner struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTaskAlert struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTaskSchedule struct {
	ID string `xml:"id,attr,omitempty"`
}

type ModifyTaskPreferences struct {
	Preference []*ModifyTaskPreferencesPreference `xml:"preference,omitempty"`
}

type ModifyTaskPreferencesPreference struct {
	ScannerName string `xml:"scanner_name,omitempty"`
	Value       string `xml:"value,omitempty"`
}

type ModifyTaskFile struct {
	Name   string `xml:"name,attr,omitempty"`
	Action string `xml:"action,attr,omitempty"`
	Data   string `xml:",chardata,omitempty"`
}

// ModifyTaskResponse represents a modify_task command response.
type ModifyTaskResponse struct {
	XMLName    xml.Name `xml:"modify_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
