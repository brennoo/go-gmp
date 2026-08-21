package filtering

import (
	"fmt"
	"strings"
)

// BuildFilter creates a filter string from a list of options.
func BuildFilter(opts ...Option) string {
	fb := NewFilter()
	for _, opt := range opts {
		opt(fb)
	}
	return fb.Build()
}

// BuildFilterWithPagination creates a filter string with pagination from options.
func BuildFilterWithPagination(page, pageSize int, opts ...Option) string {
	fb := NewFilter()

	// Add pagination first
	if page > 0 {
		fb.Offset((page-1)*pageSize + 1)
	}
	if pageSize > 0 {
		fb.Limit(pageSize)
	}

	// Add other options
	for _, opt := range opts {
		opt(fb)
	}

	return fb.Build()
}

// BuildConsolidatedFilter creates a filter string from either string filters or functional options.
// If opts are provided, they take precedence over string filters.
func BuildConsolidatedFilter(stringFilters []string, opts ...Option) string {
	// If functional options are provided, use them
	if len(opts) > 0 {
		return BuildFilter(opts...)
	}

	// Otherwise, combine string filters
	if len(stringFilters) == 0 {
		return ""
	}

	// Join string filters with spaces
	return strings.Join(stringFilters, " ")
}

// FilterArg represents either a string filter or functional options.
// This allows methods to accept both types in a single parameter.
type FilterArg interface{}

// BuildFilterFromArgs creates a filter string from mixed arguments.
// It accepts both string filters and functional options.
func BuildFilterFromArgs(args ...FilterArg) (string, error) {
	var stringFilters []string
	var opts []Option

	// Separate string filters from functional options
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			stringFilters = append(stringFilters, v)
		case Option:
			opts = append(opts, v)
		case []string:
			stringFilters = append(stringFilters, v...)
		case []Option:
			opts = append(opts, v...)
		default:
			// A value of any other type compiles (FilterArg is interface{})
			// but was never going to filter anything; erroring here is the
			// difference between a caller finding out now and a request
			// silently going out unfiltered.
			return "", fmt.Errorf("filtering: unsupported FilterArg type %T — want string, Option, []string or []Option", v)
		}
	}

	// Build filter from functional options
	var filterParts []string
	if len(opts) > 0 {
		filterParts = append(filterParts, BuildFilter(opts...))
	}

	// Add string filters
	if len(stringFilters) > 0 {
		filterParts = append(filterParts, strings.Join(stringFilters, " "))
	}

	// Combine all parts
	if len(filterParts) == 0 {
		return "", nil
	}

	return strings.Join(filterParts, " "), nil
}
