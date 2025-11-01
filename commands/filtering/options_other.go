package filtering

// WithTag adds a tag filter condition.
func WithTag(tag string) Option {
	return WithField("tag", tag)
}

// WithTagID adds a tag_id filter condition.
func WithTagID(tagID string) Option {
	return WithField("tag_id", tagID)
}

// WithType adds a type filter condition.
func WithType(typeStr string) Option {
	return WithField("type", typeStr)
}

// WithTypeLike adds a type filter condition using LIKE operator.
func WithTypeLike(pattern string) Option {
	return WithFieldLike("type", pattern)
}

// WithValue adds a value filter condition.
func WithValue(value string) Option {
	return WithField("value", value)
}

// WithValueLike adds a value filter condition using LIKE operator.
func WithValueLike(pattern string) Option {
	return WithFieldLike("value", pattern)
}

// WithUsageType adds a usage_type filter condition.
func WithUsageType(usageType string) Option {
	return WithField("usage_type", usageType)
}

// WithDescription adds a description filter condition.
func WithDescription(description string) Option {
	return WithField("description", description)
}

// WithDescriptionLike adds a description filter condition using LIKE operator.
func WithDescriptionLike(pattern string) Option {
	return WithFieldLike("description", pattern)
}

// WithAlterable adds an alterable filter condition.
func WithAlterable(alterable bool) Option {
	return WithBooleanField("alterable", alterable)
}

// WithWritable adds a writable filter condition.
func WithWritable(writable bool) Option {
	return WithBooleanField("writable", writable)
}

// WithInUse adds an in_use filter condition.
func WithInUse(inUse bool) Option {
	return WithBooleanField("in_use", inUse)
}

// WithTrash adds a trash filter condition.
func WithTrash(trash bool) Option {
	return WithBooleanField("trash", trash)
}

// WithApplyOverrides adds an apply_overrides filter condition.
func WithApplyOverrides(apply bool) Option {
	return WithBooleanField("apply_overrides", apply)
}

// WithMinQoD adds a minimum QoD filter condition.
func WithMinQoD(qod int) Option {
	return WithNumericFieldGreaterThanOrEqual("min_qod", qod)
}

// WithScanner adds a scanner filter condition.
func WithScanner(scanner string) Option {
	return WithField("scanner", scanner)
}

// WithScannerLike adds a scanner filter condition using LIKE operator.
func WithScannerLike(pattern string) Option {
	return WithFieldLike("scanner", pattern)
}

// WithAssignedTo adds an assigned_to filter condition.
func WithAssignedTo(assignedTo string) Option {
	return WithField("assigned_to", assignedTo)
}

// WithAssignedBy adds an assigned_by filter condition.
func WithAssignedBy(assignedBy string) Option {
	return WithField("assigned_by", assignedBy)
}
