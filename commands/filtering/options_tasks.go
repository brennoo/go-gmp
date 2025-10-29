package filtering

import (
	"fmt"
)

// WithTaskID adds a uuid filter condition for querying tasks by UUID.
// Note: Uses "uuid" instead of "task_id" because task_id is a direct attribute
// and doesn't work properly in filter strings. The uuid filter keyword works correctly.
func WithTaskID(taskID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("uuid", OpEqual, taskID)
	}
}

// WithTaskName adds a task name filter condition.
func WithTaskName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task_name", OpEqual, name)
	}
}

// WithTaskNameLike adds a task name filter condition using LIKE operator.
func WithTaskNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task_name", OpLike, pattern)
	}
}

// WithTaskStatus adds a task status filter condition.
func WithTaskStatus(status string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task_status", OpEqual, status)
	}
}

// WithTaskStatusIn adds a task status filter condition for multiple values.
func WithTaskStatusIn(statuses ...string) Option {
	return func(fb *FilterBuilder) {
		for _, status := range statuses {
			fb.AddCondition("task_status", OpEqual, status)
		}
	}
}

// WithTaskProgress adds a task progress filter condition.
func WithTaskProgress(progress int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("progress", OpEqual, fmt.Sprintf("%d", progress))
	}
}

// WithTaskProgressGreaterThan adds a task progress filter condition.
func WithTaskProgressGreaterThan(progress int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("progress", OpGreaterThan, fmt.Sprintf("%d", progress))
	}
}

// WithTaskProgressGreaterThanOrEqual adds a task progress filter condition.
func WithTaskProgressGreaterThanOrEqual(progress int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("progress", OpGreaterThanOrEqual, fmt.Sprintf("%d", progress))
	}
}

// WithTaskProgressLessThan adds a task progress filter condition.
func WithTaskProgressLessThan(progress int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("progress", OpLessThan, fmt.Sprintf("%d", progress))
	}
}

// WithTaskProgressLessThanOrEqual adds a task progress filter condition.
func WithTaskProgressLessThanOrEqual(progress int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("progress", OpLessThanOrEqual, fmt.Sprintf("%d", progress))
	}
}

// WithTaskAlterable adds an alterable filter condition.
func WithTaskAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, fmt.Sprintf("%t", alterable))
	}
}

// WithTaskWritable adds a writable filter condition.
func WithTaskWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, fmt.Sprintf("%t", writable))
	}
}

// WithTaskInUse adds an in_use filter condition.
func WithTaskInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, fmt.Sprintf("%t", inUse))
	}
}

// WithTaskTrash adds a trash filter condition.
func WithTaskTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, fmt.Sprintf("%t", trash))
	}
}

// WithTaskSchedule adds a schedule filter condition.
func WithTaskSchedule(schedule string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("schedule", OpEqual, schedule)
	}
}

// WithTaskScheduleLike adds a schedule filter condition using LIKE operator.
func WithTaskScheduleLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("schedule", OpLike, pattern)
	}
}

// WithTaskScanner adds a scanner filter condition.
func WithTaskScanner(scanner string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpEqual, scanner)
	}
}

// WithTaskScannerLike adds a scanner filter condition using LIKE operator.
func WithTaskScannerLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpLike, pattern)
	}
}

// WithTaskTarget adds a target filter condition.
func WithTaskTarget(target string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("target", OpEqual, target)
	}
}

// WithTaskTargetLike adds a target filter condition using LIKE operator.
func WithTaskTargetLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("target", OpLike, pattern)
	}
}

// WithTaskConfig adds a config filter condition.
func WithTaskConfig(config string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("config", OpEqual, config)
	}
}

// WithTaskConfigLike adds a config filter condition using LIKE operator.
func WithTaskConfigLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("config", OpLike, pattern)
	}
}
