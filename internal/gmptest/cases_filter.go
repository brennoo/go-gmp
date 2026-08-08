package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetFilters.Trash/Alerts are *bool: unset, true, and an explicit false
	// all round-trip distinctly.
	RegisterRequestCase("get_filters", "trash_and_alerts", &commands.GetFilters{
		Filter: "type=task",
		Trash:  commands.Bool(true),
		Alerts: commands.Bool(true),
	})

	RegisterRequestCase("get_filters", "by_id", &commands.GetFilters{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_filters", "trash_and_alerts_false", &commands.GetFilters{
		Filter: "type=task",
		Trash:  commands.Bool(false),
		Alerts: commands.Bool(false),
	})
}
