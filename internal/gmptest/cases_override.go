package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateOverride.NewSeverity/Severity are *float64: unset, a nonzero
	// value, and an explicit zero all round-trip distinctly.
	RegisterRequestCase("create_override", "unset", &commands.CreateOverride{
		Text: "False positive: patched via vendor advisory.",
		NVT:  commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
	})

	RegisterRequestCase("create_override", "severities_set", &commands.CreateOverride{
		Text:        "Downgrading severity based on compensating controls.",
		NVT:         commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		NewSeverity: commands.Float64(2.5),
		Severity:    commands.Float64(7.5),
	})

	RegisterRequestCase("create_override", "severities_zero", &commands.CreateOverride{
		Text:        "Overriding severity to zero (informational).",
		NVT:         commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		NewSeverity: commands.Float64(0),
		Severity:    commands.Float64(0),
	})

	// CreateOverride.Result/Task are ID-bearing struct pointers, not bare
	// strings, and Active is *int: this case exercises all three together.
	RegisterRequestCase("create_override", "with_result_task_active", &commands.CreateOverride{
		Text:   "Active for 7 days, scoped to one result and one task.",
		NVT:    commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		Active: commands.Int(7),
		Result: &commands.OverrideResult{ID: "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f"},
		Task:   &commands.OverrideTask{ID: "c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f"},
	})

	// ModifyOverride.NewSeverity/Severity are *float64: unset, a nonzero
	// value, and an explicit zero all round-trip distinctly.
	RegisterRequestCase("modify_override", "unset", &commands.ModifyOverride{
		OverrideID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Text:       "Updated override text.",
	})

	RegisterRequestCase("modify_override", "severities_set", &commands.ModifyOverride{
		OverrideID:  "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Text:        "Updated override text.",
		NewSeverity: commands.Float64(2.5),
		Severity:    commands.Float64(7.5),
	})

	RegisterRequestCase("modify_override", "severities_zero", &commands.ModifyOverride{
		OverrideID:  "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Text:        "Updated override text.",
		NewSeverity: commands.Float64(0),
		Severity:    commands.Float64(0),
	})

	// ModifyOverride.Result/Task are ID-bearing struct pointers, not bare
	// strings, and Active is *int: this case exercises all three together.
	RegisterRequestCase("modify_override", "with_result_task_active", &commands.ModifyOverride{
		OverrideID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
		Text:       "Active for 7 days, scoped to one result and one task.",
		Active:     commands.Int(7),
		Result:     &commands.OverrideResult{ID: "9b7f6e5d-4c3b-2a1f-0e9d-8c7b6a5d4e3f"},
		Task:       &commands.OverrideTask{ID: "c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f"},
	})
}
