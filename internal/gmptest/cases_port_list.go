package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	RegisterRequestCase("create_port_list", "basic", &commands.CreatePortList{
		Name:      "All TCP",
		PortRange: "T:1-65535",
	})

	RegisterRequestCase("create_port_list", "comment_and_copy", &commands.CreatePortList{
		Name:      "Copy of All TCP",
		PortRange: "T:1-65535",
		Comment:   "Cloned for a custom scan config",
		Copy:      "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	// GetPortLists.Details/Targets/Trash are *bool: unset, true, and an
	// explicit false all round-trip distinctly.
	RegisterRequestCase("get_port_lists", "unset", &commands.GetPortLists{})

	RegisterRequestCase("get_port_lists", "by_id", &commands.GetPortLists{
		PortListID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_port_lists", "filt_id", &commands.GetPortLists{
		FiltID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_port_lists", "all_details_true", &commands.GetPortLists{
		Filter:  "rows=-1",
		Details: commands.Bool(true),
		Targets: commands.Bool(true),
		Trash:   commands.Bool(true),
	})

	RegisterRequestCase("get_port_lists", "all_details_false", &commands.GetPortLists{
		Filter:  "rows=-1",
		Details: commands.Bool(false),
		Targets: commands.Bool(false),
		Trash:   commands.Bool(false),
	})

	// ModifyPortList.Name/Comment are plain omitempty strings, not pointers:
	// gvmd clears whichever one is left out of the request, so these cases
	// pin the set-together, name-only, and comment-only wire shapes side by
	// side rather than an omitted/explicit-empty distinction.
	RegisterRequestCase("modify_port_list", "all_unset", &commands.ModifyPortList{
		PortListID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("modify_port_list", "name_and_comment", &commands.ModifyPortList{
		PortListID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Name:       "Renamed port list",
		Comment:    "Renamed and re-commented together to avoid clearing either",
	})

	RegisterRequestCase("modify_port_list", "name_only", &commands.ModifyPortList{
		PortListID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Name:       "Renamed port list",
	})

	RegisterRequestCase("modify_port_list", "comment_only", &commands.ModifyPortList{
		PortListID: "9b16f6c2-4c9d-4e6a-8f6a-1a2b3c4d5e6f",
		Comment:    "Re-commented without resending the name",
	})
}
