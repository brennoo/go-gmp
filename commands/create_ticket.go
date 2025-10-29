package commands

import "encoding/xml"

// CreateTicket represents a create_ticket command request.
type CreateTicket struct {
	XMLName  xml.Name         `xml:"create_ticket"`
	Comment  string           `xml:"comment,omitempty"`
	Copy     string           `xml:"copy,omitempty"`
	Result   TicketResult     `xml:"result"`
	Assigned TicketAssignedTo `xml:"assigned_to"`
	OpenNote string           `xml:"open_note"`
}

// CreateTicketResponse represents a create_ticket command response.
type CreateTicketResponse struct {
	XMLName    xml.Name `xml:"create_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
	ID         string   `xml:"id,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *CreateTicketResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
