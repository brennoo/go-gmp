package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetRoles.Trash is *GMPBool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("get_roles", "filter_and_trash_true", &commands.GetRoles{
		Filter: "rows=-1",
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Trash:  commands.Bool(true),
	})

	RegisterRequestCase("get_roles", "by_id", &commands.GetRoles{
		RoleID: "b493b7a8-7489-11df-a3ec-002264764cea",
	})

	RegisterRequestCase("get_roles", "filter_and_trash_false", &commands.GetRoles{
		Filter: "rows=-1",
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Trash:  commands.Bool(false),
	})
}
