package commands

import "encoding/xml"

// MoveTask represents a move_task command request.
type MoveTask struct {
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

// GetStatus returns the status and status text from the response.
func (r *MoveTaskResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
