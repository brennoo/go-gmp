package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// CreateTag.Active is *bool: unset, true, and an explicit false all
	// round-trip distinctly. Resources is schema-required, so all three cases
	// give it the same populated value and vary only Active — a zero-value
	// Resources would emit a meaningless empty <resources><type></type>
	// </resources> that documents nothing about Active.
	tagResources := commands.TagResources{
		Type:     "task",
		Resource: []commands.TagResource{{ID: "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31"}},
	}

	RegisterRequestCase("create_tag", "active_true", &commands.CreateTag{
		Name:      "environment",
		Value:     "production",
		Resources: tagResources,
		Active:    commands.Bool(true),
	})

	RegisterRequestCase("create_tag", "active_unset", &commands.CreateTag{
		Name:      "environment",
		Value:     "staging",
		Resources: tagResources,
	})

	RegisterRequestCase("create_tag", "active_false", &commands.CreateTag{
		Name:      "environment",
		Value:     "decommissioned",
		Resources: tagResources,
		Active:    commands.Bool(false),
	})

	// ModifyTag.Active is *bool: unset, true, and an explicit false all
	// round-trip distinctly.
	RegisterRequestCase("modify_tag", "active_true", &commands.ModifyTag{
		TagID:  "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Value:  ptr("production"),
		Active: commands.Bool(true),
	})

	RegisterRequestCase("modify_tag", "active_unset", &commands.ModifyTag{
		TagID:   "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Comment: ptr("no longer actively enforced"),
	})

	RegisterRequestCase("modify_tag", "active_false", &commands.ModifyTag{
		TagID:  "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Active: commands.Bool(false),
	})

	// ModifyTag.Value/.Comment are *string (internal/codegen/overrides.go's
	// PointerOverrides): a live test against gvmd 26.36.1 proved omitting
	// <value>/<comment> leaves the stored value unchanged while sending it
	// explicitly empty clears it — Name showed no such distinction and keeps
	// the mechanical string default. See docs/review-2026-08.md's
	// post-tls_certificate revisit section. These three cases pin all three
	// wire shapes side by side.
	RegisterRequestCase("modify_tag", "value_comment_unset", &commands.ModifyTag{
		TagID: "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
	})

	RegisterRequestCase("modify_tag", "value_comment_empty", &commands.ModifyTag{
		TagID:   "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Value:   ptr(""),
		Comment: ptr(""),
	})

	RegisterRequestCase("modify_tag", "value_comment_set", &commands.ModifyTag{
		TagID:   "f1e2d3c4-b5a6-4978-8695-a4b3c2d1e0f9",
		Value:   ptr("decommissioned"),
		Comment: ptr("Rotated during Q3 audit"),
	})
}
