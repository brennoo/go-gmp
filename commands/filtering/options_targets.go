package filtering

import (
	"fmt"
)

// WithTargetID adds a target_id filter condition.
func WithTargetID(targetID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("target_id", OpEqual, targetID)
	}
}

// WithTargetName adds a target name filter condition.
func WithTargetName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("target_name", OpEqual, name)
	}
}

// WithTargetNameLike adds a target name filter condition using LIKE operator.
func WithTargetNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("target_name", OpLike, pattern)
	}
}

// WithTargetHosts adds a hosts filter condition.
func WithTargetHosts(hosts string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpEqual, hosts)
	}
}

// WithTargetHostsLike adds a hosts filter condition using LIKE operator.
func WithTargetHostsLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLike, pattern)
	}
}

// WithTargetPorts adds a ports filter condition.
func WithTargetPorts(ports string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpEqual, ports)
	}
}

// WithTargetPortsLike adds a ports filter condition using LIKE operator.
func WithTargetPortsLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpLike, pattern)
	}
}

// WithTargetPortList adds a port_list filter condition.
func WithTargetPortList(portList string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port_list", OpEqual, portList)
	}
}

// WithTargetPortListLike adds a port_list filter condition using LIKE operator.
func WithTargetPortListLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port_list", OpLike, pattern)
	}
}

// WithTargetCredentials adds a credentials filter condition.
func WithTargetCredentials(credentials string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("credentials", OpEqual, credentials)
	}
}

// WithTargetCredentialsLike adds a credentials filter condition using LIKE operator.
func WithTargetCredentialsLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("credentials", OpLike, pattern)
	}
}

// WithTargetAlterable adds an alterable filter condition.
func WithTargetAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, fmt.Sprintf("%t", alterable))
	}
}

// WithTargetWritable adds a writable filter condition.
func WithTargetWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, fmt.Sprintf("%t", writable))
	}
}

// WithTargetInUse adds an in_use filter condition.
func WithTargetInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, fmt.Sprintf("%t", inUse))
	}
}

// WithTargetTrash adds a trash filter condition.
func WithTargetTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, fmt.Sprintf("%t", trash))
	}
}

// WithTargetTasks adds a tasks filter condition.
func WithTargetTasks(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpEqual, fmt.Sprintf("%d", tasks))
	}
}

// WithTargetTasksGreaterThan adds a tasks filter condition.
func WithTargetTasksGreaterThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThan, fmt.Sprintf("%d", tasks))
	}
}

// WithTargetTasksGreaterThanOrEqual adds a tasks filter condition.
func WithTargetTasksGreaterThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThanOrEqual, fmt.Sprintf("%d", tasks))
	}
}

// WithTargetTasksLessThan adds a tasks filter condition.
func WithTargetTasksLessThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThan, fmt.Sprintf("%d", tasks))
	}
}

// WithTargetTasksLessThanOrEqual adds a tasks filter condition.
func WithTargetTasksLessThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThanOrEqual, fmt.Sprintf("%d", tasks))
	}
}

// WithTasks adds a tasks filter condition.
func WithTasks(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpEqual, fmt.Sprintf("%d", tasks))
	}
}

// WithTasksGreaterThan adds a tasks filter condition.
func WithTasksGreaterThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThan, fmt.Sprintf("%d", tasks))
	}
}

// WithTasksGreaterThanOrEqual adds a tasks filter condition.
func WithTasksGreaterThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThanOrEqual, fmt.Sprintf("%d", tasks))
	}
}

// WithTasksLessThan adds a tasks filter condition.
func WithTasksLessThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThan, fmt.Sprintf("%d", tasks))
	}
}

// WithTasksLessThanOrEqual adds a tasks filter condition.
func WithTasksLessThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThanOrEqual, fmt.Sprintf("%d", tasks))
	}
}
