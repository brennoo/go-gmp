package filtering

import (
	"fmt"
	"time"
)

// Generic helper functions for common filter patterns

// WithField adds a string field filter condition with equal operator.
func WithField(field string, value string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpEqual, value)
	}
}

// WithFieldLike adds a string field filter condition with like operator.
func WithFieldLike(field string, pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpLike, pattern)
	}
}

// WithFieldNotEqual adds a string field filter condition with not equal operator.
func WithFieldNotEqual(field string, value string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpNotEqual, value)
	}
}

// WithFieldRegex adds a string field filter condition with regex operator.
func WithFieldRegex(field string, pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpRegex, pattern)
	}
}

// WithFieldNotRegex adds a string field filter condition with not regex operator.
func WithFieldNotRegex(field string, pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpNotRegex, pattern)
	}
}

// WithNumericField adds a numeric field filter condition with equal operator.
func WithNumericField(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpEqual, fmt.Sprintf("%d", value))
	}
}

// WithNumericFieldGreaterThan adds a numeric field filter condition with greater than operator.
func WithNumericFieldGreaterThan(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpGreaterThan, fmt.Sprintf("%d", value))
	}
}

// WithNumericFieldGreaterThanOrEqual adds a numeric field filter condition with greater than or equal operator.
func WithNumericFieldGreaterThanOrEqual(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpGreaterThanOrEqual, fmt.Sprintf("%d", value))
	}
}

// WithNumericFieldLessThan adds a numeric field filter condition with less than operator.
func WithNumericFieldLessThan(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpLessThan, fmt.Sprintf("%d", value))
	}
}

// WithNumericFieldLessThanOrEqual adds a numeric field filter condition with less than or equal operator.
func WithNumericFieldLessThanOrEqual(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpLessThanOrEqual, fmt.Sprintf("%d", value))
	}
}

// WithNumericFieldNotEqual adds a numeric field filter condition with not equal operator.
func WithNumericFieldNotEqual(field string, value int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpNotEqual, fmt.Sprintf("%d", value))
	}
}

// WithBooleanField adds a boolean field filter condition with equal operator.
func WithBooleanField(field string, value bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpEqual, fmt.Sprintf("%t", value))
	}
}

// WithTimeFieldAfter adds a time field filter condition with greater than operator.
func WithTimeFieldAfter(field string, t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithTimeFieldBefore adds a time field filter condition with less than operator.
func WithTimeFieldBefore(field string, t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithTimeFieldBetween adds a time field filter condition with between operator (two conditions).
func WithTimeFieldBetween(field string, start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition(field, OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition(field, OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}

// WithFieldIn adds multiple equal conditions for a field (OR logic).
func WithFieldIn(field string, values ...string) Option {
	return func(fb *FilterBuilder) {
		for _, value := range values {
			fb.AddCondition(field, OpEqual, value)
		}
	}
}

// WithNumericFieldIn adds multiple equal conditions for a numeric field (OR logic).
func WithNumericFieldIn(field string, values ...int) Option {
	return func(fb *FilterBuilder) {
		for _, value := range values {
			fb.AddCondition(field, OpEqual, fmt.Sprintf("%d", value))
		}
	}
}
