package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateOverride.NewSeverity/Severity are *float64: unset, a nonzero
	// value, and an explicit zero all round-trip distinctly.
	RegisterRequestCase("create_override", "unset", &commands.CreateOverride{
		Text: "False positive: patched via vendor advisory.",
		NVT:  &commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
	})

	RegisterRequestCase("create_override", "severities_set", &commands.CreateOverride{
		Text:        "Downgrading severity based on compensating controls.",
		NVT:         &commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		NewSeverity: commands.Float64(2.5),
		Severity:    commands.Float64(7.5),
	})

	RegisterRequestCase("create_override", "severities_zero", &commands.CreateOverride{
		Text:        "Overriding severity to zero (informational).",
		NVT:         &commands.OverrideNVT{OID: "1.3.6.1.4.1.25623.1.0.100315"},
		NewSeverity: commands.Float64(0),
		Severity:    commands.Float64(0),
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
}
