// Package filtering builds GMP filter strings — the "first=1 rows=100
// status=Running sort=name" syntax gvmd accepts on get_* commands — from
// either plain strings or typed, composable options.
//
// Every list method on [github.com/brennoo/go-gmp/client.Client] (GetTasks,
// GetAssets, GetResults, ...) accepts a variadic [FilterArg], which is either
// a raw filter string or an [Option] built with one of the With* helpers in
// this package. [BuildFilterFromArgs] combines them into the single filter
// string gvmd expects, so the two styles can be mixed freely.
//
// A raw string is the quickest way to reuse a filter copied from the
// Greenbone web UI:
//
//	tasksResp, err := cli.GetTasks(ctx, "status=Running sort=name")
//
// The functional options give the same result without hand-building the
// syntax, and compose from named constants for status, sort fields and
// operators:
//
//	tasksResp, err := cli.GetTasks(ctx,
//		filtering.WithStatus(filtering.StatusRunning),
//		filtering.WithSortAsc(filtering.SortByName),
//		filtering.WithLimit(50),
//	)
//
// For conditions the With* helpers don't cover, build one directly with
// [NewFilter]:
//
//	filter := filtering.NewFilter().
//		AddCondition("severity", filtering.OpGreaterThanOrEqual, "7.0").
//		Sort("severity", filtering.SortDesc).
//		Build()
//	resultsResp, err := cli.GetResults(ctx, filter)
package filtering
