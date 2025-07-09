package commands

import (
	"encoding/xml"
)

// StopTaskCommand represents a stop_task command request.
type StopTaskCommand struct {
	XMLName xml.Name `xml:"stop_task"`
	TaskID  string   `xml:"task_id,attr,omitempty"`
}

// StopTaskResponse represents a stop_task command response.
type StopTaskResponse struct {
	XMLName    xml.Name `xml:"stop_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
