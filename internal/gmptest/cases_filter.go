package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetFilters.Trash/Alerts are *bool: unset, true, and an explicit false
	// all round-trip distinctly.
	RegisterRequestCase("get_filters", "trash_and_alerts", &commands.GetFilters{
		Filter: "type=task",
		Trash:  commands.Bool(true),
		Alerts: commands.Bool(true),
	})

	RegisterRequestCase("get_filters", "by_id", &commands.GetFilters{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("get_filters", "trash_and_alerts_false", &commands.GetFilters{
		Filter: "type=task",
		Trash:  commands.Bool(false),
		Alerts: commands.Bool(false),
	})

	// ModifyFilter.Name/.Comment/.Term are *string (PointerOverrides): a live
	// test against gvmd 26.36.1 proved omitting the element leaves the stored
	// value unchanged while sending it explicitly empty clears it (Term resets
	// to gvmd's default term instead of ""). See docs/review-2026-08.md's
	// `filter adopted` section. These cases pin all three wire
	// shapes per field.
	RegisterRequestCase("modify_filter", "name_comment_term_unset", &commands.ModifyFilter{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
	})

	RegisterRequestCase("modify_filter", "name_comment_term_empty", &commands.ModifyFilter{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Name:     ptr(""),
		Comment:  ptr(""),
		Term:     ptr(""),
	})

	RegisterRequestCase("modify_filter", "name_comment_term_set", &commands.ModifyFilter{
		FilterID: "3f1e9a3a-9f3a-4b3e-9c3a-1a2b3c4d5e6f",
		Name:     ptr("Renamed filter"),
		Comment:  ptr("Updated comment"),
		Term:     ptr("first=1 rows=25 sort=name"),
	})
}
