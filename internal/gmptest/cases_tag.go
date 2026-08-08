package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTag.Active is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("create_tag", "active_true", &commands.CreateTag{
		Name:  "environment",
		Value: "production",
		Resources: &commands.CreateTagResources{
			Type:      "task",
			Resources: []commands.TagResource{{ID: "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31", Type: "task"}},
		},
		Active: commands.Bool(true),
	})

	RegisterRequestCase("create_tag", "active_unset", &commands.CreateTag{
		Name:  "environment",
		Value: "staging",
	})

	RegisterRequestCase("create_tag", "active_false", &commands.CreateTag{
		Name:   "environment",
		Value:  "decommissioned",
		Active: commands.Bool(false),
	})

	// ModifyTag.Active is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("modify_tag", "active_true", &commands.ModifyTag{
		TagID:  "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Value:  "production",
		Active: commands.Bool(true),
	})

	RegisterRequestCase("modify_tag", "active_unset", &commands.ModifyTag{
		TagID:   "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Comment: "no longer actively enforced",
	})

	RegisterRequestCase("modify_tag", "active_false", &commands.ModifyTag{
		TagID:  "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Active: commands.Bool(false),
	})
}
