//nolint:dupl // Similar patterns across option types are intentional

package filtering

import (
	"fmt"
)

// WithTag adds a tag filter condition.
func WithTag(tag string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tag", OpEqual, tag)
	}
}

// WithTagID adds a tag_id filter condition.
func WithTagID(tagID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tag_id", OpEqual, tagID)
	}
}

// WithType adds a type filter condition.
func WithType(typeStr string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("type", OpEqual, typeStr)
	}
}

// WithTypeLike adds a type filter condition using LIKE operator.
func WithTypeLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("type", OpLike, pattern)
	}
}

// WithValue adds a value filter condition.
func WithValue(value string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("value", OpEqual, value)
	}
}

// WithValueLike adds a value filter condition using LIKE operator.
func WithValueLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("value", OpLike, pattern)
	}
}

// WithUsageType adds a usage_type filter condition.
func WithUsageType(usageType string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("usage_type", OpEqual, usageType)
	}
}

// WithDescription adds a description filter condition.
func WithDescription(description string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("description", OpEqual, description)
	}
}

// WithDescriptionLike adds a description filter condition using LIKE operator.
func WithDescriptionLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("description", OpLike, pattern)
	}
}

// WithAlterable adds an alterable filter condition.
func WithAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, fmt.Sprintf("%t", alterable))
	}
}

// WithWritable adds a writable filter condition.
func WithWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, fmt.Sprintf("%t", writable))
	}
}

// WithInUse adds an in_use filter condition.
func WithInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, fmt.Sprintf("%t", inUse))
	}
}

// WithTrash adds a trash filter condition.
func WithTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, fmt.Sprintf("%t", trash))
	}
}

// WithApplyOverrides adds an apply_overrides filter condition.
func WithApplyOverrides(apply bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("apply_overrides", OpEqual, fmt.Sprintf("%t", apply))
	}
}

// WithMinQoD adds a minimum QoD filter condition.
func WithMinQoD(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("min_qod", OpGreaterThanOrEqual, fmt.Sprintf("%d", qod))
	}
}

// WithScanner adds a scanner filter condition.
func WithScanner(scanner string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpEqual, scanner)
	}
}

// WithScannerLike adds a scanner filter condition using LIKE operator.
func WithScannerLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpLike, pattern)
	}
}

// WithAssignedTo adds an assigned_to filter condition.
func WithAssignedTo(assignedTo string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("assigned_to", OpEqual, assignedTo)
	}
}

// WithAssignedBy adds an assigned_by filter condition.
func WithAssignedBy(assignedBy string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("assigned_by", OpEqual, assignedBy)
	}
}
