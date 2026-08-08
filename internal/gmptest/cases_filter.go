package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetFilters.Trash/Alerts: non-pointer bools, omitempty.
	RegisterRequestCase("get_filters", "trash_and_alerts", &commands.GetFilters{
		Filter: "type=task",
		Trash:  true,
		Alerts: true,
	})

	RegisterRequestCase("get_filters", "by_id", &commands.GetFilters{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})
}
