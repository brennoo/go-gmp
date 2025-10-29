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

// GetStatus returns the status and status text from the response.
func (r *DeleteTaskResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
