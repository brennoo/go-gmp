package filtering

import (
	"testing"
	"time"
)

func TestBuildFilterFromArgs(t *testing.T) {
	// Test with no arguments
	filter := BuildFilterFromArgs()
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with string filters only
	filter = BuildFilterFromArgs("name=test", "status=running")
	expected := "name=test status=running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with functional options only
	filter = BuildFilterFromArgs(WithName("test"), WithStatus(StatusRunning))
	expected = "name=test status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with mixed arguments (both functional options and string filters are combined)
	filter = BuildFilterFromArgs("name=old", WithName("new"), "status=old", WithStatus(StatusRunning))
	expected = "name=new status=Running name=old status=old"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with string slices
	filter = BuildFilterFromArgs([]string{"name=test", "status=running"})
	expected = "name=test status=running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with option slices
	filter = BuildFilterFromArgs([]Option{WithName("test"), WithStatus(StatusRunning)})
	expected = "name=test status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}

	// Test with mixed slices
	filter = BuildFilterFromArgs([]string{"name=test"}, []Option{WithStatus(StatusRunning)})
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

func TestBuildFilterFromArgsWithInvalidTypes(t *testing.T) {
	// Test with invalid types (should be ignored)
	filter := BuildFilterFromArgs("name=test", 123, WithStatus(StatusRunning), true)
	expected := "status=Running name=test"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestBuildFilterFromArgsWithEmptySlices(t *testing.T) {
	// Test with empty string slice
	filter := BuildFilterFromArgs([]string{})
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with empty option slice
	filter = BuildFilterFromArgs([]Option{})
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}
}

func TestBuildFilterFromArgsWithNilSlices(t *testing.T) {
	// Test with nil string slice
	var stringSlice []string
	filter := BuildFilterFromArgs(stringSlice)
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}

	// Test with nil option slice
	var optionSlice []Option
	filter = BuildFilterFromArgs(optionSlice)
	if filter != "" {
		t.Errorf("Expected empty filter, got %q", filter)
	}
}

func TestBuildFilterFromArgsComplexScenario(t *testing.T) {
	// Test complex scenario with multiple types
	now := time.Now()
	filter := BuildFilterFromArgs(
		"name=test",               // String filter
		WithStatus(StatusRunning), // Functional option
		[]string{"severity>5.0"},  // String slice
		[]Option{WithLimit(10)},   // Option slice
		WithCreatedAfter(now),     // Functional option
		"sort=name",               // String filter
	)

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
