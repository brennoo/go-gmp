package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_ticket", "basic", &commands.CreateTicket{
		Result:     commands.TicketResult{ID: "5e5e5e5e-6f6f-7070-8181-929292929292"},
		AssignedTo: commands.TicketAssignedTo{User: commands.TicketUser{ID: "d290f1ee-6c54-4b01-90e6-d701748f0851"}},
		OpenNote:   "Assigning to platform team for remediation.",
	})

	// ModifyTicket.Comment/.OpenNote/.FixedNote/.ClosedNote are *string
	// (internal/codegen/overrides.go's PointerOverrides): gvmd distinguishes
	// an omitted note (existing value unchanged) from an explicit empty one
	// (clears it). Status stays a plain string because gvmd rejects an
	// explicit-empty status with a 400 (see the CommandWarning on
	// ModifyTicket.Status). These cases pin every wire shape side by side.
	RegisterRequestCase("modify_ticket", "all_unset", &commands.ModifyTicket{
		TicketID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
	})

	RegisterRequestCase("modify_ticket", "comment_empty", &commands.ModifyTicket{
		TicketID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Comment:  ptr(""),
	})

	RegisterRequestCase("modify_ticket", "comment_set", &commands.ModifyTicket{
		TicketID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Comment:  ptr("Escalating to network team."),
	})

	RegisterRequestCase("modify_ticket", "notes_empty", &commands.ModifyTicket{
		TicketID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		OpenNote:   ptr(""),
		FixedNote:  ptr(""),
		ClosedNote: ptr(""),
	})

	RegisterRequestCase("modify_ticket", "notes_set", &commands.ModifyTicket{
		TicketID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		OpenNote:   ptr("Confirmed report, investigating."),
		FixedNote:  ptr("Patched and deployed to production."),
		ClosedNote: ptr("Resolved with January update."),
	})

	RegisterRequestCase("modify_ticket", "status_set", &commands.ModifyTicket{
		TicketID: "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		Status:   "Closed",
	})

	RegisterRequestCase("modify_ticket", "assigned_to_set", &commands.ModifyTicket{
		TicketID:   "254cd3ef-bbe1-4d58-859d-21b8d0c046c6",
		AssignedTo: &commands.TicketAssignedTo{User: commands.TicketUser{ID: "33e92d3e-a379-4c46-a4cf-88c8201ab710"}},
	})

	// GetTickets.Trash is *GMPBool: unset, true, and an explicit false all
	// round-trip distinctly, matching get_roles's convention.
	RegisterRequestCase("get_tickets", "unset", &commands.GetTickets{})

	RegisterRequestCase("get_tickets", "filt_id", &commands.GetTickets{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_tickets", "trash_true", &commands.GetTickets{
		Trash: commands.Bool(true),
	})

	RegisterRequestCase("get_tickets", "trash_false", &commands.GetTickets{
		Trash: commands.Bool(false),
	})
}
