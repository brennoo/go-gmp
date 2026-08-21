package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateAlert.Active is *bool: unset, true, and an explicit false all
	// round-trip distinctly. CreateAlert.Filter is schema-required
	// (non-pointer, no omitempty), so every case gives it a populated value
	// rather than blessing a meaningless empty <filter id=""></filter>,
	// matching the convention create_tag's Resources cases already follow.
	RegisterRequestCase("create_alert", "active_true", &commands.CreateAlert{
		Name:    "Critical scan complete",
		Comment: "Notify security on-call when a critical scan finishes",
		Condition: commands.AlertCondition{
			Value: "Severity at least",
			Data:  []commands.AlertData{{Name: "severity", Value: "7.0"}},
		},
		Event: commands.AlertEvent{Value: "Task run status changed"},
		Method: commands.AlertMethod{
			Value: "Email",
			Data:  []commands.AlertData{{Name: "to_address", Value: "oncall@example.com"}},
		},
		Filter: commands.AlertFilter{ID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11"},
		Active: commands.Bool(true),
	})

	RegisterRequestCase("create_alert", "active_unset", &commands.CreateAlert{
		Name:      "Draft alert",
		Condition: commands.AlertCondition{Value: "Always"},
		Event:     commands.AlertEvent{Value: "Task run status changed"},
		Method:    commands.AlertMethod{Value: "Syslog"},
		Filter:    commands.AlertFilter{ID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11"},
	})

	RegisterRequestCase("create_alert", "active_false", &commands.CreateAlert{
		Name:      "Disabled alert",
		Condition: commands.AlertCondition{Value: "Always"},
		Event:     commands.AlertEvent{Value: "Task run status changed"},
		Method:    commands.AlertMethod{Value: "Syslog"},
		Filter:    commands.AlertFilter{ID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11"},
		Active:    commands.Bool(false),
	})

	// ModifyAlert.Name/.Comment are *string (internal/codegen/overrides.go's
	// PointerOverrides), not the plain string pointerScalars' generator-wide
	// default would produce: a live test against gvmd 26.36.1 (Greenbone
	// Community Edition) proved gvmd distinguishes an omitted <name>/<comment>
	// (existing value unchanged) from an explicit empty one (<name></name>
	// clears it), a distinction a bare string+omitempty field cannot express.
	// See docs/review-2026-08.md's `alert adopted` section. These three cases
	// pin all three wire shapes side by side so that distinction cannot
	// silently collapse back into one.
	RegisterRequestCase("modify_alert", "name_comment_unset", &commands.ModifyAlert{
		AlertID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11",
	})

	RegisterRequestCase("modify_alert", "name_comment_empty", &commands.ModifyAlert{
		AlertID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11",
		Name:    ptr(""),
		Comment: ptr(""),
	})

	RegisterRequestCase("modify_alert", "name_comment_set", &commands.ModifyAlert{
		AlertID: "6a4bd0ea-6a68-4a7c-8b3a-5f6c9d9c5c11",
		Name:    ptr("Renamed alert"),
		Comment: ptr("Updated comment"),
	})
}

// ptr returns a pointer to v. It exists because a Go string literal has no
// address of its own, and the modify_alert cases above need to express
// omitted (nil), explicit-empty (ptr("")) and non-empty pointers side by
// side.
func ptr[T any](v T) *T { return &v }
