package gmp

import "encoding/xml"

// MoveTaskCommand represents a move_task command request.
type MoveTaskCommand struct {
	XMLName xml.Name `xml:"move_task"`
	TaskID  string   `xml:"task_id,attr"`
	SlaveID string   `xml:"slave_id,attr"`
}

// MoveTaskResponse represents a move_task command response.
type MoveTaskResponse struct {
	XMLName    xml.Name `xml:"move_task_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
