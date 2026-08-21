package filtering

import (
	"testing"
	"time"
)

func TestBuildFilterFromArgs(t *testing.T) {
	// Test with no arguments
	filter, err := BuildFilterFromArgs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with string filters only
	filter, err = BuildFilterFromArgs("name=test", "status=running")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "name=test status=running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with functional options only
	filter, err = BuildFilterFromArgs(WithName("test"), WithStatus(StatusRunning))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "name=test status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with mixed arguments (both functional options and string filters are combined)
	filter, err = BuildFilterFromArgs("name=old", WithName("new"), "status=old", WithStatus(StatusRunning))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "name=new status=Running name=old status=old"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with string slices
	filter, err = BuildFilterFromArgs([]string{"name=test", "status=running"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "name=test status=running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with option slices
	filter, err = BuildFilterFromArgs([]Option{WithName("test"), WithStatus(StatusRunning)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "name=test status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with mixed slices
	filter, err = BuildFilterFromArgs([]string{"name=test"}, []Option{WithStatus(StatusRunning)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "status=Running name=test"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestBuildConsolidatedFilter(t *testing.T) {
	// Test with no arguments
	filter := BuildConsolidatedFilter([]string{})
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with string filters
	filter = BuildConsolidatedFilter([]string{"name=test", "status=running"})
	expected := "name=test status=running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with functional options
	filter = BuildConsolidatedFilter([]string{}, WithName("test"), WithStatus(StatusRunning))
	expected = "name=test status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with mixed arguments (functional options take precedence)
	filter = BuildConsolidatedFilter([]string{"name=old", "status=old"}, WithName("new"), WithStatus(StatusRunning))
	expected = "name=new status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

// TestBuildFilterFromArgsRejectsUnrecognizedTypes pins the fail-closed
// contract (T-29): any FilterArg outside string/Option/[]string/[]Option must
// error, not be silently dropped from the filter.
func TestBuildFilterFromArgsRejectsUnrecognizedTypes(t *testing.T) {
	tests := []struct {
		name string
		arg  FilterArg
	}{
		{"int", 123},
		{"bool", true},
		{"struct literal", struct{}{}},
		{"typed nil pointer", (*int)(nil)},
		// The real-world case: a caller mistakenly passes a typed request
		// struct instead of a string/Option filter (see T-28's
		// TestExplicitEmptySemantics, which did exactly this against
		// GetOverrides/GetConfigs/GetPortLists and got an unfiltered query
		// back with no error at all).
		{"unrelated struct pointer", &struct{ OverrideID string }{OverrideID: "abc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BuildFilterFromArgs("name=test", tt.arg, WithStatus(StatusRunning))
			if err == nil {
				t.Fatalf("BuildFilterFromArgs(%v) succeeded; want an unsupported-type error", tt.arg)
			}
		})
	}
}

func TestBuildFilterFromArgsWithEmptySlices(t *testing.T) {
	// Test with empty string slice
	filter, err := BuildFilterFromArgs([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with empty option slice
	filter, err = BuildFilterFromArgs([]Option{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}
}

func TestBuildFilterFromArgsWithNilSlices(t *testing.T) {
	// Test with nil string slice
	var stringSlice []string
	filter, err := BuildFilterFromArgs(stringSlice)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with nil option slice
	var optionSlice []Option
	filter, err = BuildFilterFromArgs(optionSlice)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}
}

func TestBuildFilterFromArgsComplexScenario(t *testing.T) {
	// Test complex scenario with multiple types
	now := time.Now()
	filter, err := BuildFilterFromArgs(
		"name=test",               // String filter
		WithStatus(StatusRunning), // Functional option
		[]string{"severity>5.0"},  // String slice
		[]Option{WithLimit(10)},   // Option slice
		WithCreatedAfter(now),     // Functional option
		"sort=name",               // String filter
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain all the functional options and string filters
	if !containsHelper(filter, "name=test") {
		t.Errorf("Expected filter to contain name=test, got %q", filter)
	}
	if !containsHelper(filter, "status=Running") {
		t.Errorf("Expected filter to contain status=Running, got %q", filter)
	}
	if !containsHelper(filter, "severity>5.0") {
		t.Errorf("Expected filter to contain severity>5.0, got %q", filter)
	}
	if !containsHelper(filter, "rows=10") {
		t.Errorf("Expected filter to contain rows=10, got %q", filter)
	}
	if !containsHelper(filter, "created>") {
		t.Errorf("Expected filter to contain created>, got %q", filter)
	}
	if !containsHelper(filter, "sort=name") {
		t.Errorf("Expected filter to contain sort=name, got %q", filter)
	}
}
