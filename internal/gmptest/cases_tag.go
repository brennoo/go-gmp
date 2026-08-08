package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTag.Active: non-pointer bool, omitempty.
	RegisterRequestCase("create_tag", "active_true", &commands.CreateTag{
		Name:  "environment",
		Value: "production",
		Resources: &commands.CreateTagResources{
			Type:      "task",
			Resources: []commands.TagResource{{ID: "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31", Type: "task"}},
		},
		Active: true,
	})

	RegisterRequestCase("create_tag", "active_unset", &commands.CreateTag{
		Name:  "environment",
		Value: "staging",
	})

	// ModifyTag.Active: non-pointer bool, omitempty.
	RegisterRequestCase("modify_tag", "active_true", &commands.ModifyTag{
		TagID:  "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Value:  "production",
		Active: true,
	})

	RegisterRequestCase("modify_tag", "active_unset", &commands.ModifyTag{
		TagID:   "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Comment: "no longer actively enforced",
	})
}
