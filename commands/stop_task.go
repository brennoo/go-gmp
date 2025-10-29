package commands

import (
	"encoding/xml"
)

// StopTask represents a stop_task command request.
type StopTask struct {
	XMLName xml.Name `xml:"stop_task"`
	TaskID  string   `xml:"task_id,attr,omitempty"`
}

// StopTaskResponse represents a stop_task command response.
type StopTaskResponse struct {
	XMLName    xml.Name `xml:"stop_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *StopTaskResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
