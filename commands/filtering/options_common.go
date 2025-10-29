package filtering

import (
	"strconv"
)

// WithName adds a name filter condition.
func WithName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("name", OpEqual, name)
	}
}

// WithNameLike adds a name filter condition using LIKE operator.
func WithNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("name", OpLike, pattern)
	}
}

// WithComment adds a comment filter condition.
func WithComment(comment string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("comment", OpEqual, comment)
	}
}

// WithCommentLike adds a comment filter condition using LIKE operator.
func WithCommentLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("comment", OpLike, pattern)
	}
}

// WithOwner adds an owner filter condition.
func WithOwner(owner string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("owner", OpEqual, owner)
	}
}

// WithStatus adds a status filter condition.
func WithStatus(status string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("status", OpEqual, status)
	}
}

// WithStatusIn adds a status filter condition for multiple values.
func WithStatusIn(statuses ...string) Option {
	return func(fb *FilterBuilder) {
		for _, status := range statuses {
			fb.AddCondition("status", OpEqual, status)
		}
	}
}

// WithSort adds a sort filter condition.
func WithSort(field string, direction string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(field, direction)
	}
}

// WithSortAsc adds an ascending sort filter condition.
func WithSortAsc(field string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(field, SortAsc)
	}
}

// WithSortDesc adds a descending sort filter condition.
func WithSortDesc(field string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(field, SortDesc)
	}
}

// WithLimit adds a limit filter condition.
func WithLimit(limit int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("rows", OpEqual, strconv.Itoa(limit))
	}
}

// WithOffset adds an offset filter condition.
func WithOffset(offset int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("first", OpEqual, strconv.Itoa(offset))
	}
}
