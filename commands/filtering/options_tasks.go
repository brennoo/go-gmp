package filtering

// WithTaskID adds a uuid filter condition for querying tasks by UUID.
// Note: Uses "uuid" instead of "task_id" because task_id is a direct attribute
// and doesn't work properly in filter strings. The uuid filter keyword works correctly.
func WithTaskID(taskID string) Option {
	return WithField("uuid", taskID)
}

// WithTaskName adds a task name filter condition.
func WithTaskName(name string) Option {
	return WithField("task_name", name)
}

// WithTaskNameLike adds a task name filter condition using LIKE operator.
func WithTaskNameLike(pattern string) Option {
	return WithFieldLike("task_name", pattern)
}

// WithTaskStatus adds a task status filter condition.
func WithTaskStatus(status string) Option {
	return WithField("task_status", status)
}

// WithTaskStatusIn adds a task status filter condition for multiple values.
func WithTaskStatusIn(statuses ...string) Option {
	return WithFieldIn("task_status", statuses...)
}

// WithTaskProgress adds a task progress filter condition.
func WithTaskProgress(progress int) Option {
	return WithNumericField("progress", progress)
}

// WithTaskProgressGreaterThan adds a task progress filter condition.
func WithTaskProgressGreaterThan(progress int) Option {
	return WithNumericFieldGreaterThan("progress", progress)
}

// WithTaskProgressGreaterThanOrEqual adds a task progress filter condition.
func WithTaskProgressGreaterThanOrEqual(progress int) Option {
	return WithNumericFieldGreaterThanOrEqual("progress", progress)
}

// WithTaskProgressLessThan adds a task progress filter condition.
func WithTaskProgressLessThan(progress int) Option {
	return WithNumericFieldLessThan("progress", progress)
}

// WithTaskProgressLessThanOrEqual adds a task progress filter condition.
func WithTaskProgressLessThanOrEqual(progress int) Option {
	return WithNumericFieldLessThanOrEqual("progress", progress)
}

// WithTaskAlterable adds an alterable filter condition.
func WithTaskAlterable(alterable bool) Option {
	return WithBooleanField("alterable", alterable)
}

// WithTaskWritable adds a writable filter condition.
func WithTaskWritable(writable bool) Option {
	return WithBooleanField("writable", writable)
}

// WithTaskInUse adds an in_use filter condition.
func WithTaskInUse(inUse bool) Option {
	return WithBooleanField("in_use", inUse)
}

// WithTaskTrash adds a trash filter condition.
func WithTaskTrash(trash bool) Option {
	return WithBooleanField("trash", trash)
}

// WithTaskSchedule adds a schedule filter condition.
func WithTaskSchedule(schedule string) Option {
	return WithField("schedule", schedule)
}

// WithTaskScheduleLike adds a schedule filter condition using LIKE operator.
func WithTaskScheduleLike(pattern string) Option {
	return WithFieldLike("schedule", pattern)
}

// WithTaskScanner adds a scanner filter condition.
func WithTaskScanner(scanner string) Option {
	return WithField("scanner", scanner)
}

// WithTaskScannerLike adds a scanner filter condition using LIKE operator.
func WithTaskScannerLike(pattern string) Option {
	return WithFieldLike("scanner", pattern)
}

// WithTaskTarget adds a target filter condition.
func WithTaskTarget(target string) Option {
	return WithField("target", target)
}

// WithTaskTargetLike adds a target filter condition using LIKE operator.
func WithTaskTargetLike(pattern string) Option {
	return WithFieldLike("target", pattern)
}

// WithTaskConfig adds a config filter condition.
func WithTaskConfig(config string) Option {
	return WithField("config", config)
}

// WithTaskConfigLike adds a config filter condition using LIKE operator.
func WithTaskConfigLike(pattern string) Option {
	return WithFieldLike("config", pattern)
}
