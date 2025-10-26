package filtering

import (
	"testing"
	"time"
)

func TestFilterBuilder(t *testing.T) {
	fb := NewFilter()

	// Test basic conditions
	fb.AddCondition("name", OpEqual, "test")
	fb.AddCondition("status", OpEqual, "running")

	result := fb.Build()
	expected := "name=test status=running"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFilterBuilderWithSort(t *testing.T) {
	fb := NewFilter()

	fb.AddCondition("name", OpEqual, "test")
	fb.Sort("created", "desc")

	result := fb.Build()
	expected := "name=test sort-reverse=created"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFilterBuilderWithPagination(t *testing.T) {
	fb := NewFilter()

	fb.AddCondition("name", OpEqual, "test")
	fb.Offset(10)
	fb.Limit(5)

	result := fb.Build()
	expected := "first=10 rows=5 name=test"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFunctionalOptions(t *testing.T) {
	opts := []Option{
		WithName("test-task"),
		WithStatus(StatusRunning),
		WithSeverityGreaterThan(5.0),
		WithSortAsc(SortByName),
		WithLimit(10),
	}

	filter := BuildFilter(opts...)
	expected := "name=test-task status=Running severity>5.0 rows=10 sort=name"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestFunctionalOptionsWithPagination(t *testing.T) {
	opts := []Option{
		WithName("test-task"),
		WithStatus(StatusRunning),
	}

	filter := BuildFilterWithPagination(2, 5, opts...)
	expected := "first=6 rows=5 name=test-task status=Running"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestTimeFilters(t *testing.T) {
	now := time.Now()
	opts := []Option{
		WithCreatedAfter(now),
		WithCreatedBefore(now.Add(24 * time.Hour)),
	}

	filter := BuildFilter(opts...)
	// Just check that it contains the expected parts
	if !containsHelper(filter, "created>") {
		t.Errorf("Expected filter to contain created>, got %q", filter)
	}
	if !containsHelper(filter, "created<") {
		t.Errorf("Expected filter to contain created<, got %q", filter)
	}
}

func TestBooleanFilters(t *testing.T) {
	opts := []Option{
		WithInUse(true),
		WithAlterable(false),
		WithWritable(true),
	}

	filter := BuildFilter(opts...)
	expected := "in_use=true alterable=false writable=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestIntegerFilters(t *testing.T) {
	opts := []Option{
		WithHosts(10),
		WithHostsGreaterThan(5),
		WithQoDGreaterThanOrEqual(70),
	}

	filter := BuildFilter(opts...)
	expected := "hosts=10 hosts>5 qod>=70.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestFloatFilters(t *testing.T) {
	opts := []Option{
		WithSeverity(7.5),
		WithSeverityLessThan(8.0),
	}

	filter := BuildFilter(opts...)
	expected := "severity=7.5 severity<8.0"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestEscapeValue(t *testing.T) {
	fb := NewFilter()

	// Test value with spaces
	fb.AddCondition("name", OpEqual, "test task")
	result := fb.Build()
	expected := `name="test task"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}

	// Test value with quotes
	fb2 := NewFilter()
	fb2.AddCondition("name", OpEqual, `test"task`)
	result2 := fb2.Build()
	expected2 := `name="test\"task"`
	if result2 != expected2 {
		t.Errorf("Expected %q, got %q", expected2, result2)
	}
}
