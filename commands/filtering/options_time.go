//nolint:dupl // Similar patterns across option types are intentional

package filtering

import (
	"time"
)

// WithCreatedAfter adds a creation time filter condition.
func WithCreatedAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("created", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithCreatedBefore adds a creation time filter condition.
func WithCreatedBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("created", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithCreatedBetween adds a creation time filter condition.
func WithCreatedBetween(start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("created", OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition("created", OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}

// WithModifiedAfter adds a modification time filter condition.
func WithModifiedAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("modified", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithModifiedBefore adds a modification time filter condition.
func WithModifiedBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("modified", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithModifiedBetween adds a modification time filter condition.
func WithModifiedBetween(start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("modified", OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition("modified", OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}

// WithNextDueAfter adds a next due time filter condition.
func WithNextDueAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("next_due", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithNextDueBefore adds a next due time filter condition.
func WithNextDueBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("next_due", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithNextDueBetween adds a next due time filter condition.
func WithNextDueBetween(start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("next_due", OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition("next_due", OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}

// WithFirstReportCreatedAfter adds a first_report_created filter condition.
func WithFirstReportCreatedAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("first_report_created", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithFirstReportCreatedBefore adds a first_report_created filter condition.
func WithFirstReportCreatedBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("first_report_created", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithFirstReportCreatedBetween adds a first_report_created filter condition.
func WithFirstReportCreatedBetween(start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("first_report_created", OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition("first_report_created", OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}

// WithLastReportCreatedAfter adds a last_report_created filter condition.
func WithLastReportCreatedAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("last_report_created", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithLastReportCreatedBefore adds a last_report_created filter condition.
func WithLastReportCreatedBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("last_report_created", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithLastReportCreatedBetween adds a last_report_created filter condition.
func WithLastReportCreatedBetween(start, end time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("last_report_created", OpGreaterThanOrEqual, start.Format("2006-01-02T15:04:05Z"))
		fb.AddCondition("last_report_created", OpLessThanOrEqual, end.Format("2006-01-02T15:04:05Z"))
	}
}
