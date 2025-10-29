package commands

import "encoding/xml"

// ModifyTicket represents a modify_ticket command request.
type ModifyTicket struct {
	XMLName    xml.Name          `xml:"modify_ticket"`
	TicketID   string            `xml:"ticket_id,attr"`
	Comment    string            `xml:"comment,omitempty"`
	Status     string            `xml:"status,omitempty"`
	OpenNote   string            `xml:"open_note,omitempty"`
	FixedNote  string            `xml:"fixed_note,omitempty"`
	ClosedNote string            `xml:"closed_note,omitempty"`
	Assigned   *TicketAssignedTo `xml:"assigned_to,omitempty"`
}

// ModifyTicketResponse represents a modify_ticket command response.
type ModifyTicketResponse struct {
	XMLName    xml.Name `xml:"modify_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *ModifyTicketResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
