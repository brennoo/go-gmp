package filtering

import (
	"fmt"
	"strings"
	"time"
)

// Filter represents a filter that can be applied to GMP commands.
type Filter interface {
	// Build returns the filter string for the GMP protocol.
	Build() string
}

// FilterBuilder provides a fluent interface for building filters.
type FilterBuilder struct {
	conditions []string
	sort       string
	limit      int
	offset     int
}

// NewFilter creates a new filter builder.
func NewFilter() *FilterBuilder {
	return &FilterBuilder{
		conditions: make([]string, 0),
	}
}

// AddCondition adds a filter condition.
func (fb *FilterBuilder) AddCondition(column, operator, value string) *FilterBuilder {
	// Escape value if it contains spaces or special characters
	escapedValue := escapeValue(value)
	fb.conditions = append(fb.conditions, fmt.Sprintf("%s%s%s", column, operator, escapedValue))
	return fb
}

// AddConditionf adds a filter condition with formatting.
func (fb *FilterBuilder) AddConditionf(column, operator, format string, args ...interface{}) *FilterBuilder {
	value := fmt.Sprintf(format, args...)
	return fb.AddCondition(column, operator, value)
}

// Sort sets the sort order.
func (fb *FilterBuilder) Sort(column, order string) *FilterBuilder {
	if order == "desc" || order == "descending" {
		fb.sort = fmt.Sprintf("sort-reverse=%s", column)
	} else {
		fb.sort = fmt.Sprintf("sort=%s", column)
	}
	return fb
}

// Limit sets the maximum number of results.
func (fb *FilterBuilder) Limit(n int) *FilterBuilder {
	fb.limit = n
	return fb
}

// Offset sets the starting position (1-based).
func (fb *FilterBuilder) Offset(n int) *FilterBuilder {
	fb.offset = n
	return fb
}

// Build constructs the final filter string.
func (fb *FilterBuilder) Build() string {
	parts := make([]string, 0)

	// Add pagination
	if fb.offset > 0 {
		parts = append(parts, fmt.Sprintf("first=%d", fb.offset))
	}
	if fb.limit > 0 {
		parts = append(parts, fmt.Sprintf("rows=%d", fb.limit))
	}

	// Add conditions
	parts = append(parts, fb.conditions...)

	// Add sorting
	if fb.sort != "" {
		parts = append(parts, fb.sort)
	}

	return strings.Join(parts, " ")
}

// StringFilter represents a string-based filter condition.
type StringFilter struct {
	column   string
	operator string
	value    string
}

// NewStringFilter creates a new string filter.
func NewStringFilter(column, operator, value string) *StringFilter {
	return &StringFilter{
		column:   column,
		operator: operator,
		value:    value,
	}
}

// Build returns the filter string.
func (sf *StringFilter) Build() string {
	return fmt.Sprintf("%s%s%s", sf.column, sf.operator, escapeValue(sf.value))
}

// IntFilter represents an integer-based filter condition.
type IntFilter struct {
	column   string
	operator string
	value    int
}

// NewIntFilter creates a new integer filter.
func NewIntFilter(column, operator string, value int) *IntFilter {
	return &IntFilter{
		column:   column,
		operator: operator,
		value:    value,
	}
}

// Build returns the filter string.
func (if_ *IntFilter) Build() string {
	return fmt.Sprintf("%s%s%d", if_.column, if_.operator, if_.value)
}

// FloatFilter represents a float-based filter condition.
type FloatFilter struct {
	column   string
	operator string
	value    float64
}

// NewFloatFilter creates a new float filter.
func NewFloatFilter(column, operator string, value float64) *FloatFilter {
	return &FloatFilter{
		column:   column,
		operator: operator,
		value:    value,
	}
}

// Build returns the filter string.
func (ff *FloatFilter) Build() string {
	return fmt.Sprintf("%s%s%.1f", ff.column, ff.operator, ff.value)
}

// TimeFilter represents a time-based filter condition.
type TimeFilter struct {
	column   string
	operator string
	value    time.Time
}

// NewTimeFilter creates a new time filter.
func NewTimeFilter(column, operator string, value time.Time) *TimeFilter {
	return &TimeFilter{
		column:   column,
		operator: operator,
		value:    value,
	}
}

// Build returns the filter string.
func (tf *TimeFilter) Build() string {
	// Format time as ISO 8601
	isoTime := tf.value.Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf("%s%s%s", tf.column, tf.operator, escapeValue(isoTime))
}

// CompositeFilter represents multiple filter conditions combined.
type CompositeFilter struct {
	filters []Filter
	sort    string
	limit   int
	offset  int
}

// NewCompositeFilter creates a new composite filter.
func NewCompositeFilter(filters ...Filter) *CompositeFilter {
	return &CompositeFilter{
		filters: filters,
	}
}

// AddFilter adds a filter to the composite.
func (cf *CompositeFilter) AddFilter(filter Filter) *CompositeFilter {
	cf.filters = append(cf.filters, filter)
	return cf
}

// Sort sets the sort order.
func (cf *CompositeFilter) Sort(column, order string) *CompositeFilter {
	if order == "desc" || order == "descending" {
		cf.sort = fmt.Sprintf("sort-reverse=%s", column)
	} else {
		cf.sort = fmt.Sprintf("sort=%s", column)
	}
	return cf
}

// Limit sets the maximum number of results.
func (cf *CompositeFilter) Limit(n int) *CompositeFilter {
	cf.limit = n
	return cf
}

// Offset sets the starting position (1-based).
func (cf *CompositeFilter) Offset(n int) *CompositeFilter {
	cf.offset = n
	return cf
}

// Build constructs the final filter string.
func (cf *CompositeFilter) Build() string {
	parts := make([]string, 0)

	// Add pagination
	if cf.offset > 0 {
		parts = append(parts, fmt.Sprintf("first=%d", cf.offset))
	}
	if cf.limit > 0 {
		parts = append(parts, fmt.Sprintf("rows=%d", cf.limit))
	}

	// Add filter conditions
	for _, filter := range cf.filters {
		if filter != nil {
			parts = append(parts, filter.Build())
		}
	}

	// Add sorting
	if cf.sort != "" {
		parts = append(parts, cf.sort)
	}

	return strings.Join(parts, " ")
}

// escapeValue escapes a value for use in filter strings.
func escapeValue(value string) string {
	// If value contains spaces or special characters, wrap in quotes
	if strings.ContainsAny(value, " \t\n\r\"'=<>!") {
		return fmt.Sprintf(`"%s"`, strings.ReplaceAll(value, `"`, `\"`))
	}
	return value
}
