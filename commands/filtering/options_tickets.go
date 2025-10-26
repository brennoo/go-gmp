//nolint:dupl // Similar patterns across option types are intentional

package filtering

import (
	"fmt"
)

// WithTicketID adds a ticket_id filter condition.
func WithTicketID(ticketID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ticket_id", OpEqual, ticketID)
	}
}

// WithTicketName adds a ticket name filter condition.
func WithTicketName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ticket_name", OpEqual, name)
	}
}

// WithTicketNameLike adds a ticket name filter condition using LIKE operator.
func WithTicketNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ticket_name", OpLike, pattern)
	}
}

// WithTicketStatus adds a ticket status filter condition.
func WithTicketStatus(status string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ticket_status", OpEqual, status)
	}
}

// WithTicketStatusIn adds a ticket status filter condition for multiple values.
func WithTicketStatusIn(statuses ...string) Option {
	return func(fb *FilterBuilder) {
		for _, status := range statuses {
			fb.AddCondition("ticket_status", OpEqual, status)
		}
	}
}

// WithTicketAssignedTo adds an assigned_to filter condition.
func WithTicketAssignedTo(assignedTo string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("assigned_to", OpEqual, assignedTo)
	}
}

// WithTicketAssignedBy adds an assigned_by filter condition.
func WithTicketAssignedBy(assignedBy string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("assigned_by", OpEqual, assignedBy)
	}
}

// WithTicketResult adds a result filter condition.
func WithTicketResult(result string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result", OpEqual, result)
	}
}

// WithTicketResultLike adds a result filter condition using LIKE operator.
func WithTicketResultLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result", OpLike, pattern)
	}
}

// WithTicketTask adds a task filter condition.
func WithTicketTask(task string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpEqual, task)
	}
}

// WithTicketTaskLike adds a task filter condition using LIKE operator.
func WithTicketTaskLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpLike, pattern)
	}
}

// WithTicketHost adds a host filter condition.
func WithTicketHost(host string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpEqual, host)
	}
}

// WithTicketHostLike adds a host filter condition using LIKE operator.
func WithTicketHostLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpLike, pattern)
	}
}

// WithTicketPort adds a port filter condition.
func WithTicketPort(port string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port", OpEqual, port)
	}
}

// WithTicketPortLike adds a port filter condition using LIKE operator.
func WithTicketPortLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port", OpLike, pattern)
	}
}

// WithTicketNVT adds an NVT filter condition.
func WithTicketNVT(nvt string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("nvt", OpEqual, nvt)
	}
}

// WithTicketNVTLike adds an NVT filter condition using LIKE operator.
func WithTicketNVTLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("nvt", OpLike, pattern)
	}
}

// WithTicketAlterable adds an alterable filter condition.
func WithTicketAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, fmt.Sprintf("%t", alterable))
	}
}

// WithTicketWritable adds a writable filter condition.
func WithTicketWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, fmt.Sprintf("%t", writable))
	}
}

// WithTicketInUse adds an in_use filter condition.
func WithTicketInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, fmt.Sprintf("%t", inUse))
	}
}

// WithTicketTrash adds a trash filter condition.
func WithTicketTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, fmt.Sprintf("%t", trash))
	}
}
