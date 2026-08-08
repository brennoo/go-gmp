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

	// GetAggregates.FirstGroup/MaxGroups are *int: unset (the case above),
	// a nonzero value, and an explicit zero all round-trip distinctly.
	RegisterRequestCase("get_aggregates", "first_and_max_groups_set", &commands.GetAggregates{
		Type:        "task",
		GroupColumn: "severity",
		FirstGroup:  commands.Int(1),
		MaxGroups:   commands.Int(10),
	})

	RegisterRequestCase("get_aggregates", "first_and_max_groups_zero", &commands.GetAggregates{
		Type:        "task",
		GroupColumn: "severity",
		FirstGroup:  commands.Int(0),
		MaxGroups:   commands.Int(0),
	})

	// GetSettings.First/Max are *int: unset, a nonzero value, and an
	// explicit zero all round-trip distinctly.
	RegisterRequestCase("get_settings", "unset", &commands.GetSettings{
		Filter: "rows=-1",
	})

	RegisterRequestCase("get_settings", "first_and_max_set", &commands.GetSettings{
		Filter: "rows=-1",
		First:  commands.Int(1),
		Max:    commands.Int(50),
	})

	RegisterRequestCase("get_settings", "first_and_max_zero", &commands.GetSettings{
		Filter: "rows=-1",
		First:  commands.Int(0),
		Max:    commands.Int(0),
	})
}
