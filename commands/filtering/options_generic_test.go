package filtering

import (
	"testing"
	"time"
)

func TestWithField(t *testing.T) {
	opts := []Option{WithField("test_field", "test_value")}
	filter := BuildFilter(opts...)
	expected := "test_field=test_value"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldLike(t *testing.T) {
	opts := []Option{WithFieldLike("test_field", "test%")}
	filter := BuildFilter(opts...)
	expected := "test_field~test%"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldNotEqual(t *testing.T) {
	opts := []Option{WithFieldNotEqual("test_field", "test_value")}
	filter := BuildFilter(opts...)
	expected := "test_field!=test_value"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldRegex(t *testing.T) {
	opts := []Option{WithFieldRegex("test_field", "test.*")}
	filter := BuildFilter(opts...)
	expected := "test_field=~test.*"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldNotRegex(t *testing.T) {
	opts := []Option{WithFieldNotRegex("test_field", "test.*")}
	filter := BuildFilter(opts...)
	expected := "test_field!~test.*"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericField(t *testing.T) {
	opts := []Option{WithNumericField("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field=42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldGreaterThan(t *testing.T) {
	opts := []Option{WithNumericFieldGreaterThan("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field>42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldGreaterThanOrEqual(t *testing.T) {
	opts := []Option{WithNumericFieldGreaterThanOrEqual("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field>=42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldLessThan(t *testing.T) {
	opts := []Option{WithNumericFieldLessThan("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field<42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldLessThanOrEqual(t *testing.T) {
	opts := []Option{WithNumericFieldLessThanOrEqual("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field<=42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldNotEqual(t *testing.T) {
	opts := []Option{WithNumericFieldNotEqual("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field!=42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithBooleanField(t *testing.T) {
	opts := []Option{WithBooleanField("test_field", true)}
	filter := BuildFilter(opts...)
	expected := "test_field=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithBooleanFieldFalse(t *testing.T) {
	opts := []Option{WithBooleanField("test_field", false)}
	filter := BuildFilter(opts...)
	expected := "test_field=false"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTimeFieldAfter(t *testing.T) {
	testTime := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	opts := []Option{WithTimeFieldAfter("test_field", testTime)}
	filter := BuildFilter(opts...)
	expected := "test_field>2023-01-15T10:30:00Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTimeFieldBefore(t *testing.T) {
	testTime := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	opts := []Option{WithTimeFieldBefore("test_field", testTime)}
	filter := BuildFilter(opts...)
	expected := "test_field<2023-01-15T10:30:00Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithTimeFieldBetween(t *testing.T) {
	startTime := time.Date(2023, 1, 15, 10, 30, 0, 0, time.UTC)
	endTime := time.Date(2023, 1, 16, 10, 30, 0, 0, time.UTC)
	opts := []Option{WithTimeFieldBetween("test_field", startTime, endTime)}
	filter := BuildFilter(opts...)
	expected := "test_field>=2023-01-15T10:30:00Z test_field<=2023-01-16T10:30:00Z"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldIn(t *testing.T) {
	opts := []Option{WithFieldIn("test_field", "value1", "value2", "value3")}
	filter := BuildFilter(opts...)
	expected := "test_field=value1 test_field=value2 test_field=value3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldIn(t *testing.T) {
	opts := []Option{WithNumericFieldIn("test_field", 1, 2, 3)}
	filter := BuildFilter(opts...)
	expected := "test_field=1 test_field=2 test_field=3"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithFieldInSingleValue(t *testing.T) {
	opts := []Option{WithFieldIn("test_field", "single_value")}
	filter := BuildFilter(opts...)
	expected := "test_field=single_value"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestWithNumericFieldInSingleValue(t *testing.T) {
	opts := []Option{WithNumericFieldIn("test_field", 42)}
	filter := BuildFilter(opts...)
	expected := "test_field=42"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestMultipleGenericFilters(t *testing.T) {
	opts := []Option{
		WithField("name", "test"),
		WithNumericField("count", 5),
		WithBooleanField("active", true),
	}
	filter := BuildFilter(opts...)
	expected := "name=test count=5 active=true"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}

func TestGenericFiltersWithExistingOptions(t *testing.T) {
	opts := []Option{
		WithField("name", "test"),
		WithSort("created", "desc"),
		WithLimit(10),
	}
	filter := BuildFilter(opts...)
	expected := "name=test rows=10 sort-reverse=created"
	if filter != expected {
		t.Errorf("Expected %q, got %q", expected, filter)
	}
}
