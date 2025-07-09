package commands

import (
	"encoding/xml"
)

// GetTasks represents a get_tasks command request.
type GetTasks struct {
	XMLName          xml.Name `xml:"get_tasks"`
	TaskID           string   `xml:"task_id,attr,omitempty"`
	Filter           string   `xml:"filter,attr,omitempty"`
	FiltID           string   `xml:"filt_id,attr,omitempty"`
	Trash            bool     `xml:"trash,attr,omitempty"`
	Details          bool     `xml:"details,attr,omitempty"`
	IgnorePagination bool     `xml:"ignore_pagination,attr,omitempty"`
	SchedulesOnly    bool     `xml:"schedules_only,attr,omitempty"`
	UsageType        string   `xml:"usage_type,attr,omitempty"`
}

// GetTasksResponse represents a get_tasks command response.
type GetTasksResponse struct {
	XMLName        xml.Name                  `xml:"get_tasks_response"`
	Status         string                    `xml:"status,attr"`
	StatusText     string                    `xml:"status_text,attr"`
	ApplyOverrides string                    `xml:"apply_overrides"`
	Task           []GetTasksResponseTask    `xml:"task"`
	Filters        GetTasksResponseFilters   `xml:"filters"`
	Sort           GetTasksResponseSort      `xml:"sort"`
	Tasks          GetTasksResponseTasks     `xml:"tasks"`
	TaskCount      GetTasksResponseTaskCount `xml:"task_count"`
}

// GetTasksResponseTask represents a task element in the get_tasks response.
type GetTasksResponseTask struct {
	ID               string           `xml:"id,attr"`
	Owner            *TaskOwner       `xml:"owner,omitempty"`
	Name             string           `xml:"name"`
	Comment          string           `xml:"comment,omitempty"`
	CreationTime     string           `xml:"creation_time,omitempty"`
	ModificationTime string           `xml:"modification_time,omitempty"`
	Writable         string           `xml:"writable,omitempty"`
	InUse            string           `xml:"in_use,omitempty"`
	Permissions      *TaskPermissions `xml:"permissions,omitempty"`
	Alterable        string           `xml:"alterable,omitempty"`
	UsageType        string           `xml:"usage_type,omitempty"`
	Config           *TaskConfig      `xml:"config,omitempty"`
	Target           *TaskTarget      `xml:"target,omitempty"`
	HostsOrdering    string           `xml:"hosts_ordering,omitempty"`
	Scanner          *TaskScanner     `xml:"scanner,omitempty"`
	Alert            *TaskAlert       `xml:"alert,omitempty"`
	Observers        *TaskObservers   `xml:"observers,omitempty"`
	Schedule         *TaskSchedule    `xml:"schedule,omitempty"`
	ReportCount      string           `xml:"report_count,omitempty"`
	Trend            string           `xml:"trend,omitempty"`
	Status           string           `xml:"status,omitempty"`
	Progress         string           `xml:"progress,omitempty"`
	LastReport       *TaskLastReport  `xml:"last_report,omitempty"`
	Reports          *TaskReports     `xml:"reports,omitempty"`
	Preferences      *TaskPreferences `xml:"preferences,omitempty"`
}

type TaskOwner struct {
	Name string `xml:"name"`
}

type TaskPermissions struct {
	Permission []TaskPermission `xml:"permission"`
}

type TaskPermission struct {
	Name string `xml:"name"`
}

type TaskConfig struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Trash string `xml:"trash,omitempty"`
}

type TaskTarget struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash string `xml:"trash,omitempty"`
}

type TaskScanner struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Type  string `xml:"type"`
	Trash string `xml:"trash,omitempty"`
}

type TaskAlert struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type TaskObservers struct {
	UserList string          `xml:"user_list,omitempty"`
	Groups   []ObserverGroup `xml:"group,omitempty"`
	Roles    []ObserverRole  `xml:"role,omitempty"`
}

type ObserverGroup struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type ObserverRole struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name"`
}

type TaskSchedule struct {
	ID    string `xml:"id,attr"`
	Name  string `xml:"name"`
	Trash string `xml:"trash,omitempty"`
}

type GetTasksResponseFilters struct {
	ID       string                          `xml:"id,attr"`
	Term     string                          `xml:"term"`
	Name     string                          `xml:"name"`
	Keywords GetTasksResponseFiltersKeywords `xml:"keywords"`
}

type GetTasksResponseFiltersKeywords struct {
	Keyword []GetTasksResponseFiltersKeywordsKeyword `xml:"keyword"`
}

type GetTasksResponseFiltersKeywordsKeyword struct {
	Column   string `xml:"column"`
	Relation string `xml:"relation"`
	Value    string `xml:"value"`
}

type GetTasksResponseSort struct {
	Value string                    `xml:",chardata"`
	Field GetTasksResponseSortField `xml:"field"`
}

type GetTasksResponseSortField struct {
	Order string `xml:"order"`
}

type GetTasksResponseTasks struct {
	Start int `xml:"start,attr"`
	Max   int `xml:"max,attr"`
}

type GetTasksResponseTaskCount struct {
	Filtered int `xml:"filtered"`
	Page     int `xml:"page"`
}

type TaskLastReport struct {
	Report *TaskReport `xml:"report"`
}

type TaskReports struct {
	Report []TaskReport `xml:"report"`
}

type TaskPreferences struct {
	Preference []TaskPreference `xml:"preference"`
}

type TaskPreference struct {
	Name        string `xml:"name"`
	ScannerName string `xml:"scanner_name"`
	Value       string `xml:"value"`
}

type TaskReport struct {
	ID            string `xml:"id,attr"`
	Timestamp     string `xml:"timestamp,omitempty"`
	ScanStart     string `xml:"scan_start,omitempty"`
	ScanEnd       string `xml:"scan_end,omitempty"`
	ScanRunStatus string `xml:"scan_run_status,omitempty"`
	Severity      string `xml:"severity,omitempty"`
}
