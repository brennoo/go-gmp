package commands

import "encoding/xml"

// StartTask represents a start_task command request.
type StartTask struct {
	XMLName xml.Name `xml:"start_task"`
	TaskID  string   `xml:"task_id,attr,omitempty"`
}

// StartTaskResponse represents a start_task command response.
type StartTaskResponse struct {
	XMLName    xml.Name `xml:"start_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ReportID   string   `xml:"report_id"`
}
