package filtering

import (
	"testing"
	"time"
)

func TestStringFilter(t *testing.T) {
	sf := NewStringFilter("name", OpEqual, "test")

	result := sf.Build()
	expected := "name=test"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestStringFilterWithSpaces(t *testing.T) {
	sf := NewStringFilter("name", OpEqual, "test task")

	result := sf.Build()
	expected := `name="test task"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestIntFilter(t *testing.T) {
	inf := NewIntFilter("count", OpEqual, 10)

	result := inf.Build()
	expected := "count=10"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestIntFilterWithGreaterThan(t *testing.T) {
	inf := NewIntFilter("limit", OpGreaterThan, 5)

	result := inf.Build()
	expected := "limit>5"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestIntFilterWithLessThan(t *testing.T) {
	inf := NewIntFilter("offset", OpLessThan, 100)

	result := inf.Build()
	expected := "offset<100"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFloatFilter(t *testing.T) {
	ff := NewFloatFilter("severity", OpEqual, 7.5)

	result := ff.Build()
	expected := "severity=7.5"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFloatFilterWithGreaterThan(t *testing.T) {
	ff := NewFloatFilter("score", OpGreaterThan, 5.0)

	result := ff.Build()
	expected := "score>5.0"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestFloatFilterWithLessThan(t *testing.T) {
	ff := NewFloatFilter("rating", OpLessThan, 10.0)

	result := ff.Build()
	expected := "rating<10.0"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTimeFilter(t *testing.T) {
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	tf := NewTimeFilter("created", OpEqual, now)

	result := tf.Build()
	expected := "created=2023-01-01T12:00:00Z"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTimeFilterWithGreaterThan(t *testing.T) {
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	tf := NewTimeFilter("modified", OpGreaterThan, now)

	result := tf.Build()
	expected := "modified>2023-01-01T12:00:00Z"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTimeFilterWithLessThan(t *testing.T) {
	now := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	tf := NewTimeFilter("expires", OpLessThan, now)

	result := tf.Build()
	expected := "expires<2023-01-01T12:00:00Z"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestCompositeFilter(t *testing.T) {
	cf := NewCompositeFilter()

	// Add string filter
	stringFilter := NewStringFilter("name", OpEqual, "test")
	cf.AddFilter(stringFilter)

	// Add int filter
	intFilter := NewIntFilter("count", OpGreaterThan, 5)
	cf.AddFilter(intFilter)

	// Add float filter
	floatFilter := NewFloatFilter("severity", OpLessThan, 8.0)
	cf.AddFilter(floatFilter)

	result := cf.Build()
	if !containsHelper(result, "name=test") {
		t.Errorf("Expected filter to contain name=test, got %q", result)
	}
	if !containsHelper(result, "count>5") {
		t.Errorf("Expected filter to contain count>5, got %q", result)
	}
	if !containsHelper(result, "severity<8.0") {
		t.Errorf("Expected filter to contain severity<8.0, got %q", result)
	}
}

func TestCompositeFilterWithSort(t *testing.T) {
	cf := NewCompositeFilter()

	stringFilter := NewStringFilter("name", OpEqual, "test")
	cf.AddFilter(stringFilter)

	cf.Sort("name", "desc")

	result := cf.Build()
	if !containsHelper(result, "name=test") {
		t.Errorf("Expected filter to contain name=test, got %q", result)
	}
	if !containsHelper(result, "sort-reverse=name") {
		t.Errorf("Expected filter to contain sort-reverse=name, got %q", result)
	}
}

func TestCompositeFilterWithPagination(t *testing.T) {
	cf := NewCompositeFilter()

	stringFilter := NewStringFilter("name", OpEqual, "test")
	cf.AddFilter(stringFilter)

	cf.Offset(40)
	cf.Limit(60)

	result := cf.Build()
	if !containsHelper(result, "name=test") {
		t.Errorf("Expected filter to contain name=test, got %q", result)
	}
	if !containsHelper(result, "first=40") {
		t.Errorf("Expected filter to contain first=40, got %q", result)
	}
	if !containsHelper(result, "rows=60") {
		t.Errorf("Expected filter to contain rows=60, got %q", result)
	}
}

func TestAddConditionf(t *testing.T) {
	fb := NewFilter()

	// Test with format string
	fb.AddConditionf("name", OpEqual, "test-%s", "task")

	result := fb.Build()
	expected := "name=test-task"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestAddConditionfWithMultipleArgs(t *testing.T) {
	fb := NewFilter()

	// Test with multiple format arguments
	fb.AddConditionf("pattern", OpEqual, "%s-%d-%s", "test", 123, "end")

	result := fb.Build()
	expected := "pattern=test-123-end"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestAddConditionfWithSpecialChars(t *testing.T) {
	fb := NewFilter()

	// Test with special characters that need escaping
	fb.AddConditionf("path", OpEqual, "/path/to/%s", "file with spaces")

	result := fb.Build()
	expected := `path="/path/to/file with spaces"`
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}
