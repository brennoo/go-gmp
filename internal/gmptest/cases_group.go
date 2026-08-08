package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetGroups.Trash: non-pointer bool, omitempty.
	RegisterRequestCase("get_groups", "trash_true", &commands.GetGroups{
		Filter: "rows=-1",
		Trash:  true,
	})

	RegisterRequestCase("get_groups", "by_id", &commands.GetGroups{
		GroupID: "0a1b2c3d-4e5f-6071-8293-a4b5c6d7e8f9",
	})
}
