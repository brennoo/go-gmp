package filtering

// WithTargetID adds a target_id filter condition.
func WithTargetID(targetID string) Option {
	return WithField("target_id", targetID)
}

// WithTargetName adds a target name filter condition.
func WithTargetName(name string) Option {
	return WithField("target_name", name)
}

// WithTargetNameLike adds a target name filter condition using LIKE operator.
func WithTargetNameLike(pattern string) Option {
	return WithFieldLike("target_name", pattern)
}

// WithTargetHosts adds a hosts filter condition.
func WithTargetHosts(hosts string) Option {
	return WithField("hosts", hosts)
}

// WithTargetHostsLike adds a hosts filter condition using LIKE operator.
func WithTargetHostsLike(pattern string) Option {
	return WithFieldLike("hosts", pattern)
}

// WithTargetPorts adds a ports filter condition.
func WithTargetPorts(ports string) Option {
	return WithField("ports", ports)
}

// WithTargetPortsLike adds a ports filter condition using LIKE operator.
func WithTargetPortsLike(pattern string) Option {
	return WithFieldLike("ports", pattern)
}

// WithTargetPortList adds a port_list filter condition.
func WithTargetPortList(portList string) Option {
	return WithField("port_list", portList)
}

// WithTargetPortListLike adds a port_list filter condition using LIKE operator.
func WithTargetPortListLike(pattern string) Option {
	return WithFieldLike("port_list", pattern)
}

// WithTargetCredentials adds a credentials filter condition.
func WithTargetCredentials(credentials string) Option {
	return WithField("credentials", credentials)
}

// WithTargetCredentialsLike adds a credentials filter condition using LIKE operator.
func WithTargetCredentialsLike(pattern string) Option {
	return WithFieldLike("credentials", pattern)
}

// WithTargetAlterable adds an alterable filter condition.
func WithTargetAlterable(alterable bool) Option {
	return WithBooleanField("alterable", alterable)
}

// WithTargetWritable adds a writable filter condition.
func WithTargetWritable(writable bool) Option {
	return WithBooleanField("writable", writable)
}

// WithTargetInUse adds an in_use filter condition.
func WithTargetInUse(inUse bool) Option {
	return WithBooleanField("in_use", inUse)
}

// WithTargetTrash adds a trash filter condition.
func WithTargetTrash(trash bool) Option {
	return WithBooleanField("trash", trash)
}

// WithTargetTasks adds a tasks filter condition.
func WithTargetTasks(tasks int) Option {
	return WithNumericField("tasks", tasks)
}

// WithTargetTasksGreaterThan adds a tasks filter condition.
func WithTargetTasksGreaterThan(tasks int) Option {
	return WithNumericFieldGreaterThan("tasks", tasks)
}

// WithTargetTasksGreaterThanOrEqual adds a tasks filter condition.
func WithTargetTasksGreaterThanOrEqual(tasks int) Option {
	return WithNumericFieldGreaterThanOrEqual("tasks", tasks)
}

// WithTargetTasksLessThan adds a tasks filter condition.
func WithTargetTasksLessThan(tasks int) Option {
	return WithNumericFieldLessThan("tasks", tasks)
}

// WithTargetTasksLessThanOrEqual adds a tasks filter condition.
func WithTargetTasksLessThanOrEqual(tasks int) Option {
	return WithNumericFieldLessThanOrEqual("tasks", tasks)
}

// WithTasks adds a tasks filter condition.
func WithTasks(tasks int) Option {
	return WithNumericField("tasks", tasks)
}

// WithTasksGreaterThan adds a tasks filter condition.
func WithTasksGreaterThan(tasks int) Option {
	return WithNumericFieldGreaterThan("tasks", tasks)
}

// WithTasksGreaterThanOrEqual adds a tasks filter condition.
func WithTasksGreaterThanOrEqual(tasks int) Option {
	return WithNumericFieldGreaterThanOrEqual("tasks", tasks)
}

// WithTasksLessThan adds a tasks filter condition.
func WithTasksLessThan(tasks int) Option {
	return WithNumericFieldLessThan("tasks", tasks)
}

// WithTasksLessThanOrEqual adds a tasks filter condition.
func WithTasksLessThanOrEqual(tasks int) Option {
	return WithNumericFieldLessThanOrEqual("tasks", tasks)
}
