package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTarget.ReverseLookupOnly/ReverseLookupUnify/AllowSimultaneousIPs:
	// non-pointer bools, omitempty.
	RegisterRequestCase("create_target", "flags_true", &commands.CreateTarget{
		Name:                 "Corporate LAN",
		Hosts:                "192.0.2.0/24",
		ReverseLookupOnly:    true,
		ReverseLookupUnify:   true,
		AllowSimultaneousIPs: true,
	})

	RegisterRequestCase("create_target", "flags_unset", &commands.CreateTarget{
		Name:  "Corporate LAN",
		Hosts: "192.0.2.0/24",
	})

	// GetTargets.Trash/Tasks/Details: non-pointer bools, omitempty.
	RegisterRequestCase("get_targets", "all_details", &commands.GetTargets{
		Filter:  "rows=-1",
		Trash:   true,
		Tasks:   true,
		Details: true,
	})

	RegisterRequestCase("get_targets", "by_id", &commands.GetTargets{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
	})

	// DeleteTarget.Ultimate: non-pointer bool, omitempty attr.
	RegisterRequestCase("delete_target", "ultimate_true", &commands.DeleteTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		Ultimate: true,
	})

	RegisterRequestCase("delete_target", "ultimate_unset", &commands.DeleteTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
	})

	// ModifyTarget.ReverseLookupOnly/ReverseLookupUnify/AllowSimultaneousIPs:
	// non-pointer bools, omitempty.
	RegisterRequestCase("modify_target", "flags_true", &commands.ModifyTarget{
		TargetID:             "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		ReverseLookupOnly:    true,
		ReverseLookupUnify:   true,
		AllowSimultaneousIPs: true,
	})

	RegisterRequestCase("modify_target", "flags_unset", &commands.ModifyTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		Comment:  "Renamed after the office move",
	})
}
