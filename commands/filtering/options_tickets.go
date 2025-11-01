package filtering

// WithTicketID adds a ticket_id filter condition.
func WithTicketID(ticketID string) Option {
	return WithField("ticket_id", ticketID)
}

// WithTicketName adds a ticket name filter condition.
func WithTicketName(name string) Option {
	return WithField("ticket_name", name)
}

// WithTicketNameLike adds a ticket name filter condition using LIKE operator.
func WithTicketNameLike(pattern string) Option {
	return WithFieldLike("ticket_name", pattern)
}

// WithTicketStatus adds a ticket status filter condition.
func WithTicketStatus(status string) Option {
	return WithField("ticket_status", status)
}

// WithTicketStatusIn adds a ticket status filter condition for multiple values.
func WithTicketStatusIn(statuses ...string) Option {
	return WithFieldIn("ticket_status", statuses...)
}

// WithTicketAssignedTo adds an assigned_to filter condition.
func WithTicketAssignedTo(assignedTo string) Option {
	return WithField("assigned_to", assignedTo)
}

// WithTicketAssignedBy adds an assigned_by filter condition.
func WithTicketAssignedBy(assignedBy string) Option {
	return WithField("assigned_by", assignedBy)
}

// WithTicketResult adds a result filter condition.
func WithTicketResult(result string) Option {
	return WithField("result", result)
}

// WithTicketResultLike adds a result filter condition using LIKE operator.
func WithTicketResultLike(pattern string) Option {
	return WithFieldLike("result", pattern)
}

// WithTicketTask adds a task filter condition.
func WithTicketTask(task string) Option {
	return WithField("task", task)
}

// WithTicketTaskLike adds a task filter condition using LIKE operator.
func WithTicketTaskLike(pattern string) Option {
	return WithFieldLike("task", pattern)
}

// WithTicketHost adds a host filter condition.
func WithTicketHost(host string) Option {
	return WithField("host", host)
}

// WithTicketHostLike adds a host filter condition using LIKE operator.
func WithTicketHostLike(pattern string) Option {
	return WithFieldLike("host", pattern)
}

// WithTicketPort adds a port filter condition.
func WithTicketPort(port string) Option {
	return WithField("port", port)
}

// WithTicketPortLike adds a port filter condition using LIKE operator.
func WithTicketPortLike(pattern string) Option {
	return WithFieldLike("port", pattern)
}

// WithTicketNVT adds an NVT filter condition.
func WithTicketNVT(nvt string) Option {
	return WithField("nvt", nvt)
}

// WithTicketNVTLike adds an NVT filter condition using LIKE operator.
func WithTicketNVTLike(pattern string) Option {
	return WithFieldLike("nvt", pattern)
}

// WithTicketAlterable adds an alterable filter condition.
func WithTicketAlterable(alterable bool) Option {
	return WithBooleanField("alterable", alterable)
}

// WithTicketWritable adds a writable filter condition.
func WithTicketWritable(writable bool) Option {
	return WithBooleanField("writable", writable)
}

// WithTicketInUse adds an in_use filter condition.
func WithTicketInUse(inUse bool) Option {
	return WithBooleanField("in_use", inUse)
}

// WithTicketTrash adds a trash filter condition.
func WithTicketTrash(trash bool) Option {
	return WithBooleanField("trash", trash)
}
