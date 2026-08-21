package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetGroups.Trash is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_groups", "trash_true", &commands.GetGroups{
		Filter: "rows=-1",
		Trash:  commands.Bool(true),
	})

	RegisterRequestCase("get_groups", "by_id", &commands.GetGroups{
		GroupID: "0a1b2c3d-4e5f-6071-8293-a4b5c6d7e8f9",
	})

	RegisterRequestCase("get_groups", "trash_false", &commands.GetGroups{
		Filter: "rows=-1",
		Trash:  commands.Bool(false),
	})
}
