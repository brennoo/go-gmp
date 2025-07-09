package commands

import "encoding/xml"

// CreateTask represents a create_task command request.
type CreateTask struct {
	XMLName         xml.Name               `xml:"create_task"`
	Name            string                 `xml:"name,omitempty"`
	Comment         string                 `xml:"comment,omitempty"`
	Copy            string                 `xml:"copy,omitempty"`
	Alterable       bool                   `xml:"alterable,omitempty"`
	UsageType       string                 `xml:"usage_type,omitempty"`
	Config          *CreateTaskConfig      `xml:"config,omitempty"`
	Target          *CreateTaskTarget      `xml:"target,omitempty"`
	HostsOrdering   string                 `xml:"hosts_ordering,omitempty"`
	Scanner         *CreateTaskScanner     `xml:"scanner,omitempty"`
	Alert           []*CreateTaskAlert     `xml:"alert,omitempty"`
	Schedule        *CreateTaskSchedule    `xml:"schedule,omitempty"`
	SchedulePeriods int                    `xml:"schedule_periods,omitempty"`
	Observers       *CreateTaskObservers   `xml:"observers,omitempty"`
	Preferences     *CreateTaskPreferences `xml:"preferences,omitempty"`
}

// CreateTaskConfig represents the config element for create_task.
type CreateTaskConfig struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskTarget represents the target element for create_task.
type CreateTaskTarget struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskScanner represents the scanner element for create_task.
type CreateTaskScanner struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskAlert represents the alert element for create_task.
type CreateTaskAlert struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskSchedule represents the schedule element for create_task.
type CreateTaskSchedule struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskObservers represents the observers element for create_task.
type CreateTaskObservers struct {
	Users  string                      `xml:",chardata,omitempty"`
	Groups []*CreateTaskObserversGroup `xml:"group,omitempty"`
}

// CreateTaskObserversGroup represents a group in the observers element for create_task.
type CreateTaskObserversGroup struct {
	ID string `xml:"id,attr,omitempty"`
}

// CreateTaskPreferences represents the preferences element for create_task.
type CreateTaskPreferences struct {
	Preference []*CreateTaskPreferencesPreference `xml:"preference,omitempty"`
}

// CreateTaskPreferencesPreference represents a preference in the preferences element for create_task.
type CreateTaskPreferencesPreference struct {
	ScannerName string `xml:"scanner_name,omitempty"`
	Value       string `xml:"value,omitempty"`
}

// CreateTaskResponse represents a create_task command response.
type CreateTaskResponse struct {
	XMLName    xml.Name `xml:"create_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}
