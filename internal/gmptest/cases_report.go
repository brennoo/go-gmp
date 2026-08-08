package gmptest

import "github.com/brennoo/go-gmp/commands"

func init() {
	// GetResults.NotesDetails/OverridesDetails/Details/GetCounts are *bool:
	// unset, true, and an explicit false all round-trip distinctly.
	RegisterRequestCase("get_results", "all_details", &commands.GetResults{
		Filter:           "rows=-1",
		TaskID:           "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31",
		NotesDetails:     commands.Bool(true),
		OverridesDetails: commands.Bool(true),
		Details:          commands.Bool(true),
		GetCounts:        commands.Bool(true),
	})

	RegisterRequestCase("get_results", "by_id", &commands.GetResults{
		ResultID: "5e5e5e5e-6f6f-7070-8181-929292929292",
	})

	RegisterRequestCase("get_results", "all_details_false", &commands.GetResults{
		Filter:           "rows=-1",
		TaskID:           "e512e2ca-9d0e-4bf3-bc73-7fbe6e9bbf31",
		NotesDetails:     commands.Bool(false),
		OverridesDetails: commands.Bool(false),
		Details:          commands.Bool(false),
		GetCounts:        commands.Bool(false),
	})

	RegisterRequestCase("get_aggregates", "task_severity", &commands.GetAggregates{
		Type:        "task",
		GroupColumn: "severity",
		DataColumn:  "high",
		SortField:   "severity",
		SortOrder:   "descending",
	})
}
