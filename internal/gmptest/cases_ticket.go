package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_ticket", "basic", &commands.CreateTicket{
		Result:   commands.TicketResult{ID: "5e5e5e5e-6f6f-7070-8181-929292929292"},
		Assigned: commands.TicketAssignedTo{User: &commands.TicketUser{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851"}},
		OpenNote: "Assigning to platform team for remediation.",
	})
}
