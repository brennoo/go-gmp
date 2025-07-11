package commands

import (
	"encoding/xml"
)

// DeleteTask represents a delete_task command request.
type DeleteTask struct {
	XMLName  xml.Name `xml:"delete_task"`
	TaskID   string   `xml:"task_id,attr,omitempty"`
	Ultimate bool     `xml:"ultimate,attr,omitempty"`
}

// DeleteTaskResponse represents a delete_task command response.
type DeleteTaskResponse struct {
	XMLName    xml.Name `xml:"delete_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
