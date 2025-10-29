package commands

import (
	"encoding/xml"
	"time"
)

// GetTickets represents a get_tickets command request.
type GetTickets struct {
	XMLName  xml.Name `xml:"get_tickets"`
	TicketID string   `xml:"ticket_id,attr,omitempty"`
	Filter   string   `xml:"filter,attr,omitempty"`
}

// GetTicketsResponse represents a get_tickets command response.
type GetTicketsResponse struct {
	XMLName     xml.Name       `xml:"get_tickets_response"`
	Status      string         `xml:"status,attr"`
	StatusText  string         `xml:"status_text,attr"`
	Tickets     []Ticket       `xml:"ticket"`
	Filters     *TicketFilters `xml:"filters,omitempty"`
	Sort        *TicketSort    `xml:"sort,omitempty"`
	TicketsInfo *TicketTickets `xml:"tickets,omitempty"`
	TicketCount *TicketCount   `xml:"ticket_count,omitempty"`
	Truncated   string         `xml:"truncated,omitempty"`
}

// GetStatus returns the status and status text from the response.
func (r *GetTicketsResponse) GetStatus() (string, string) {
	return r.Status, r.StatusText
}

// Ticket represents a <ticket> element in the get_tickets response.
type Ticket struct {
	ID                string             `xml:"id,attr"`
	Owner             *TicketOwner       `xml:"owner,omitempty"`
	Name              string             `xml:"name"`
	Comment           string             `xml:"comment"`
	CreationTime      time.Time          `xml:"creation_time"`
	ModificationTime  time.Time          `xml:"modification_time"`
	Writable          bool               `xml:"writable"`
	InUse             bool               `xml:"in_use"`
	Permissions       *TicketPermissions `xml:"permissions,omitempty"`
	UserTags          *TicketUserTags    `xml:"user_tags,omitempty"`
	NVT               *TicketNVT         `xml:"nvt,omitempty"`
	Result            *TicketResult      `xml:"result,omitempty"`
	AssignedTo        *TicketAssignedTo  `xml:"assigned_to,omitempty"`
	Task              *TicketTask        `xml:"task,omitempty"`
	Report            *TicketReport      `xml:"report,omitempty"`
	Severity          float64            `xml:"severity,omitempty"`
	Host              string             `xml:"host,omitempty"`
	Location          string             `xml:"location,omitempty"`
	SolutionType      string             `xml:"solution_type,omitempty"`
	Status            string             `xml:"status"`
	OpenTime          time.Time          `xml:"open_time,omitempty"`
	OpenNote          string             `xml:"open_note,omitempty"`
	FixedTime         time.Time          `xml:"fixed_time,omitempty"`
	FixedNote         string             `xml:"fixed_note,omitempty"`
	ClosedTime        time.Time          `xml:"closed_time,omitempty"`
	ClosedNote        string             `xml:"closed_note,omitempty"`
	FixVerifiedTime   time.Time          `xml:"fix_verified_time,omitempty"`
	FixVerifiedReport *TicketReport      `xml:"fix_verified_report,omitempty"`
	Orphan            bool               `xml:"orphan,omitempty"`
}

// TicketOwner represents the <owner> element for a ticket.
type TicketOwner struct {
	Name string `xml:"name"`
}

// TicketPermissions represents the <permissions> element for a ticket.
type TicketPermissions struct {
	Permission []TicketPermission `xml:"permission"`
}

// TicketPermission represents a <permission> element for a ticket.
type TicketPermission struct {
	Name string `xml:"name"`
}

// TicketAssignedTo represents the <assigned_to> element for a ticket.
type TicketAssignedTo struct {
	User *TicketUser `xml:"user"`
}

// TicketUser represents the <user> element in assigned_to.
type TicketUser struct {
	ID string `xml:"id,attr"`
}

// TicketUserTags represents the <user_tags> element for a ticket.
type TicketUserTags struct {
	Count int             `xml:"count,attr,omitempty"`
	Tag   []TicketUserTag `xml:"tag,omitempty"`
}

// TicketUserTag represents a <tag> element in user_tags.
type TicketUserTag struct {
	ID    string `xml:"id,attr,omitempty"`
	Name  string `xml:"name,omitempty"`
	Value string `xml:"value,omitempty"`
}

// TicketNVT represents the <nvt> element for a ticket.
type TicketNVT struct {
	OID  string `xml:"oid,attr"`
	Name string `xml:"name,omitempty"`
}

// TicketTask represents the <task> element for a ticket.
type TicketTask struct {
	ID string `xml:"id,attr"`
}

// TicketReport represents the <report> element for a ticket.
type TicketReport struct {
	ID string `xml:"id,attr"`
}

// TicketResult represents the <result> element for a ticket.
type TicketResult struct {
	ID string `xml:"id,attr"`
}

// TicketFilters represents the <filters> element in the response.
type TicketFilters struct {
	ID      string `xml:"id,attr,omitempty"`
	Term    string `xml:"term,omitempty"`
	Name    string `xml:"name,omitempty"`
	Comment string `xml:"comment,omitempty"`
}

// TicketSort represents the <sort> element in the response.
type TicketSort struct {
	Field string `xml:"field,omitempty"`
	Order string `xml:"order,omitempty"`
}

// TicketTickets represents the <tickets> meta element in the response.
type TicketTickets struct {
	Start int `xml:"start,attr,omitempty"`
	Max   int `xml:"max,attr,omitempty"`
}

// TicketCount represents the <ticket_count> element in the response.
type TicketCount struct {
	Filtered int `xml:"filtered,attr,omitempty"`
	Page     int `xml:"page,attr,omitempty"`
}
