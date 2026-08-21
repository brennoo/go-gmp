package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTarget.ReverseLookupOnly/ReverseLookupUnify/AllowSimultaneousIPs
	// are *bool: unset, true, and an explicit false all round-trip distinctly.
	RegisterRequestCase("create_target", "flags_true", &commands.CreateTarget{
		Name:                 "Corporate LAN",
		Hosts:                "192.0.2.0/24",
		ReverseLookupOnly:    commands.Bool(true),
		ReverseLookupUnify:   commands.Bool(true),
		AllowSimultaneousIPs: commands.Bool(true),
	})

	RegisterRequestCase("create_target", "flags_unset", &commands.CreateTarget{
		Name:  "Corporate LAN",
		Hosts: "192.0.2.0/24",
	})

	RegisterRequestCase("create_target", "flags_false", &commands.CreateTarget{
		Name:                 "Corporate LAN",
		Hosts:                "192.0.2.0/24",
		ReverseLookupOnly:    commands.Bool(false),
		ReverseLookupUnify:   commands.Bool(false),
		AllowSimultaneousIPs: commands.Bool(false),
	})

	// GetTargets.Trash/Tasks/Details are *bool: unset, true, and an explicit
	// false all round-trip distinctly.
	RegisterRequestCase("get_targets", "all_details", &commands.GetTargets{
		Filter:  "rows=-1",
		Trash:   commands.Bool(true),
		Tasks:   commands.Bool(true),
		Details: commands.Bool(true),
	})

	RegisterRequestCase("get_targets", "by_id", &commands.GetTargets{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
	})

	RegisterRequestCase("get_targets", "all_details_false", &commands.GetTargets{
		Filter:  "rows=-1",
		Trash:   commands.Bool(false),
		Tasks:   commands.Bool(false),
		Details: commands.Bool(false),
	})

	// DeleteTarget.Ultimate is a required GMPBool: true and false round-trip
	// to distinct wire values ("1"/"0"), and the attribute is always present.
	RegisterRequestCase("delete_target", "ultimate_true", &commands.DeleteTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		Ultimate: true,
	})

	RegisterRequestCase("delete_target", "ultimate_false", &commands.DeleteTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		Ultimate: false,
	})

	// ModifyTarget.ReverseLookupOnly/ReverseLookupUnify/AllowSimultaneousIPs
	// are *bool: unset, true, and an explicit false all round-trip distinctly.
	RegisterRequestCase("modify_target", "flags_true", &commands.ModifyTarget{
		TargetID:             "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		ReverseLookupOnly:    commands.Bool(true),
		ReverseLookupUnify:   commands.Bool(true),
		AllowSimultaneousIPs: commands.Bool(true),
	})

	RegisterRequestCase("modify_target", "flags_unset", &commands.ModifyTarget{
		TargetID: "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		Comment:  ptr("Renamed after the office move"),
	})

	RegisterRequestCase("modify_target", "flags_false", &commands.ModifyTarget{
		TargetID:             "1a2b3c4d-5e6f-7081-92a3-b4c5d6e7f8a9",
		ReverseLookupOnly:    commands.Bool(false),
		ReverseLookupUnify:   commands.Bool(false),
		AllowSimultaneousIPs: commands.Bool(false),
	})
}
