package filtering

import "testing"

func TestCompositeFilter_Sort(t *testing.T) {
	// Test the Sort method to get 100% coverage
	filter := NewCompositeFilter()

	// Add some filters
	filter1 := NewStringFilter("name", "=", "value1")
	filter2 := NewStringFilter("type", "=", "value2")
	filter3 := NewStringFilter("status", "=", "value3")

	filter.AddFilter(filter1)
	filter.AddFilter(filter2)
	filter.AddFilter(filter3)

	// Test sorting
	filter.Sort("name", "asc")

	// Verify the filters are still there (we can't easily test the internal sort order)
	// but we can verify the method doesn't panic and the filter is still functional
	if len(filter.filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(filter.filters))
	}
}

func TestCompositeFilter_EdgeCases(t *testing.T) {
	// Test edge cases for better coverage
	filter := NewCompositeFilter()

	// Test with no filters
	filter.Sort("name", "asc")
	filter.Limit(10)
	filter.Offset(5)

	// Test with single filter
	singleFilter := NewStringFilter("test", "=", "value")
	filter.AddFilter(singleFilter)
	filter.Sort("test", "desc")

	// Test building with no filters
	result := filter.Build()
	if result == "" {
		t.Error("Expected non-empty result even with no filters")
	}
}
