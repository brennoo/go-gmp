package commands

import "encoding/xml"

// DeleteTicket represents a delete_ticket command request.
type DeleteTicket struct {
	XMLName  xml.Name `xml:"delete_ticket"`
	TicketID string   `xml:"ticket_id,attr"`
	Ultimate bool     `xml:"ultimate,attr"`
}

// DeleteTicketResponse represents a delete_ticket command response.
type DeleteTicketResponse struct {
	XMLName    xml.Name `xml:"delete_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}

// GetStatus returns the status and status text from the response.
func (r *DeleteTicketResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}
