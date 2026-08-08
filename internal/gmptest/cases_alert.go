package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateAlert.Active is a non-pointer bool tagged omitempty: a caller who
	// leaves it false gets no <active> element at all, so this pair pins both
	// sides of that behaviour on disk.
	RegisterRequestCase("create_alert", "active_true", &commands.CreateAlert{
		Name:    "Critical scan complete",
		Comment: "Notify security on-call when a critical scan finishes",
		Condition: commands.AlertCondition{
			Text: "Severity at least",
			Data: []commands.AlertData{{Name: "severity", Text: "7.0"}},
		},
		Event: commands.AlertEvent{Text: "Task run status changed"},
		Method: commands.AlertMethod{
			Text: "Email",
			Data: []commands.AlertData{{Name: "to_address", Text: "oncall@example.com"}},
		},
		Filter: &commands.AlertFilter{ID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11"},
		Active: true,
	})

	RegisterRequestCase("create_alert", "active_unset", &commands.CreateAlert{
		Name:      "Draft alert",
		Condition: commands.AlertCondition{Text: "Always"},
		Event:     commands.AlertEvent{Text: "Task run status changed"},
		Method:    commands.AlertMethod{Text: "Syslog"},
	})
}
