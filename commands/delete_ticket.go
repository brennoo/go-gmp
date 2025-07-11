package commands

import "encoding/xml"

// DeleteTicket represents a delete_ticket command request.
type DeleteTicket struct {
	XMLName  xml.Name `xml:"delete_ticket"`
	TicketID string   `xml:"ticket_id,attr"`
	Ultimate string   `xml:"ultimate,attr"`
}

// DeleteTicketResponse represents a delete_ticket command response.
type DeleteTicketResponse struct {
	XMLName    xml.Name `xml:"delete_ticket_response"`
	Status     string   `xml:"status,attr"`
	StatusText string   `xml:"status_text,attr"`
}
