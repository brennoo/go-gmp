//nolint:dupl // Similar patterns across option types are intentional

package filtering

import (
	"fmt"
	"strconv"
	"time"
)

// Option represents a functional option that can be applied to filters.
type Option func(*FilterBuilder)

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

// WithSeverity adds a severity filter condition.
func WithSeverity(severity float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("severity", OpEqual, fmt.Sprintf("%.1f", severity))
	}
}

// WithSeverityGreaterThan adds a severity filter condition.
func WithSeverityGreaterThan(severity float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("severity", OpGreaterThan, fmt.Sprintf("%.1f", severity))
	}
}

// WithSeverityGreaterThanOrEqual adds a severity filter condition.
func WithSeverityGreaterThanOrEqual(severity float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("severity", OpGreaterThanOrEqual, fmt.Sprintf("%.1f", severity))
	}
}

// WithSeverityLessThan adds a severity filter condition.
func WithSeverityLessThan(severity float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("severity", OpLessThan, fmt.Sprintf("%.1f", severity))
	}
}

// WithSeverityLessThanOrEqual adds a severity filter condition.
func WithSeverityLessThanOrEqual(severity float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("severity", OpLessThanOrEqual, fmt.Sprintf("%.1f", severity))
	}
}

// WithThreat adds a threat filter condition.
func WithThreat(threat string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("threat", OpEqual, threat)
	}
}

// WithTrend adds a trend filter condition.
func WithTrend(trend string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trend", OpEqual, trend)
	}
}

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

// WithSort adds a sort condition.
func WithSort(column, order string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(column, order)
	}
}

// WithSortAsc adds an ascending sort condition.
func WithSortAsc(column string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(column, SortAsc)
	}
}

// WithSortDesc adds a descending sort condition.
func WithSortDesc(column string) Option {
	return func(fb *FilterBuilder) {
		fb.Sort(column, SortDesc)
	}
}

// WithLimit adds a limit condition.
func WithLimit(n int) Option {
	return func(fb *FilterBuilder) {
		fb.Limit(n)
	}
}

// WithOffset adds an offset condition.
func WithOffset(n int) Option {
	return func(fb *FilterBuilder) {
		fb.Offset(n)
	}
}

// WithTag adds a tag filter condition.
func WithTag(tag string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tag", OpEqual, tag)
	}
}

// WithTagID adds a tag ID filter condition.
func WithTagID(tagID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tag_id", OpEqual, tagID)
	}
}

// WithUsageType adds a usage type filter condition.
func WithUsageType(usageType string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("usage_type", OpEqual, usageType)
	}
}

// WithInUse adds an in_use filter condition.
func WithInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, strconv.FormatBool(inUse))
	}
}

// WithAlterable adds an alterable filter condition.
func WithAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, strconv.FormatBool(alterable))
	}
}

// WithWritable adds a writable filter condition.
func WithWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, strconv.FormatBool(writable))
	}
}

// WithTrash adds a trash filter condition.
func WithTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, strconv.FormatBool(trash))
	}
}

// WithApplyOverrides adds an apply_overrides filter condition.
func WithApplyOverrides(apply bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("apply_overrides", OpEqual, strconv.FormatBool(apply))
	}
}

// WithMinQoD adds a minimum QoD filter condition.
func WithMinQoD(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("min_qod", OpGreaterThanOrEqual, strconv.Itoa(qod))
	}
}

// WithHost adds a host filter condition.
func WithHost(host string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpEqual, host)
	}
}

// WithHostLike adds a host filter condition using LIKE operator.
func WithHostLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host", OpLike, pattern)
	}
}

// WithPort adds a port filter condition.
func WithPort(port string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port", OpEqual, port)
	}
}

// WithPortLike adds a port filter condition using LIKE operator.
func WithPortLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("port", OpLike, pattern)
	}
}

// WithNVT adds an NVT filter condition.
func WithNVT(nvt string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("nvt", OpEqual, nvt)
	}
}

// WithNVTLike adds an NVT filter condition using LIKE operator.
func WithNVTLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("nvt", OpLike, pattern)
	}
}

// WithTask adds a task filter condition.
func WithTask(task string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpEqual, task)
	}
}

// WithTaskLike adds a task filter condition using LIKE operator.
func WithTaskLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpLike, pattern)
	}
}

// WithReport adds a report filter condition.
func WithReport(report string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report", OpEqual, report)
	}
}

// WithReportLike adds a report filter condition using LIKE operator.
func WithReportLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report", OpLike, pattern)
	}
}

// WithLocation adds a location filter condition.
func WithLocation(location string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpEqual, location)
	}
}

// WithLocationLike adds a location filter condition using LIKE operator.
func WithLocationLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("location", OpLike, pattern)
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

// WithQoD adds a QoD filter condition.
func WithQoD(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpEqual, strconv.Itoa(qod))
	}
}

// WithQoDGreaterThan adds a QoD filter condition.
func WithQoDGreaterThan(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpGreaterThan, strconv.Itoa(qod))
	}
}

// WithQoDGreaterThanOrEqual adds a QoD filter condition.
func WithQoDGreaterThanOrEqual(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpGreaterThanOrEqual, strconv.Itoa(qod))
	}
}

// WithQoDLessThan adds a QoD filter condition.
func WithQoDLessThan(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpLessThan, strconv.Itoa(qod))
	}
}

// WithQoDLessThanOrEqual adds a QoD filter condition.
func WithQoDLessThanOrEqual(qod int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpLessThanOrEqual, strconv.Itoa(qod))
	}
}

// WithHosts adds a hosts filter condition.
func WithHosts(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpEqual, strconv.Itoa(hosts))
	}
}

// WithHostsGreaterThan adds a hosts filter condition.
func WithHostsGreaterThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThan, strconv.Itoa(hosts))
	}
}

// WithHostsGreaterThanOrEqual adds a hosts filter condition.
func WithHostsGreaterThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThanOrEqual, strconv.Itoa(hosts))
	}
}

// WithHostsLessThan adds a hosts filter condition.
func WithHostsLessThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThan, strconv.Itoa(hosts))
	}
}

// WithHostsLessThanOrEqual adds a hosts filter condition.
func WithHostsLessThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThanOrEqual, strconv.Itoa(hosts))
	}
}

// WithOS adds an OS filter condition.
func WithOS(os string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpEqual, os)
	}
}

// WithOSLike adds an OS filter condition using LIKE operator.
func WithOSLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("os", OpLike, pattern)
	}
}

// WithHostCount adds a host count filter condition.
func WithHostCount(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpEqual, strconv.Itoa(count))
	}
}

// WithHostCountGreaterThan adds a host count filter condition.
func WithHostCountGreaterThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThan, strconv.Itoa(count))
	}
}

// WithHostCountGreaterThanOrEqual adds a host count filter condition.
func WithHostCountGreaterThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpGreaterThanOrEqual, strconv.Itoa(count))
	}
}

// WithHostCountLessThan adds a host count filter condition.
func WithHostCountLessThan(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThan, strconv.Itoa(count))
	}
}

// WithHostCountLessThanOrEqual adds a host count filter condition.
func WithHostCountLessThanOrEqual(count int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("host_count", OpLessThanOrEqual, strconv.Itoa(count))
	}
}

// WithMaxHosts adds a max hosts filter condition.
func WithMaxHosts(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpEqual, strconv.Itoa(maxHosts))
	}
}

// WithMaxHostsGreaterThan adds a max hosts filter condition.
func WithMaxHostsGreaterThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThan, strconv.Itoa(maxHosts))
	}
}

// WithMaxHostsGreaterThanOrEqual adds a max hosts filter condition.
func WithMaxHostsGreaterThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpGreaterThanOrEqual, strconv.Itoa(maxHosts))
	}
}

// WithMaxHostsLessThan adds a max hosts filter condition.
func WithMaxHostsLessThan(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThan, strconv.Itoa(maxHosts))
	}
}

// WithMaxHostsLessThanOrEqual adds a max hosts filter condition.
func WithMaxHostsLessThanOrEqual(maxHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("max_hosts", OpLessThanOrEqual, strconv.Itoa(maxHosts))
	}
}

// WithTasks adds a tasks filter condition.
func WithTasks(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpEqual, strconv.Itoa(tasks))
	}
}

// WithTasksGreaterThan adds a tasks filter condition.
func WithTasksGreaterThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThan, strconv.Itoa(tasks))
	}
}

// WithTasksGreaterThanOrEqual adds a tasks filter condition.
func WithTasksGreaterThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpGreaterThanOrEqual, strconv.Itoa(tasks))
	}
}

// WithTasksLessThan adds a tasks filter condition.
func WithTasksLessThan(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThan, strconv.Itoa(tasks))
	}
}

// WithTasksLessThanOrEqual adds a tasks filter condition.
func WithTasksLessThanOrEqual(tasks int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("tasks", OpLessThanOrEqual, strconv.Itoa(tasks))
	}
}

// WithPorts adds a ports filter condition.
func WithPorts(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpEqual, strconv.Itoa(ports))
	}
}

// WithPortsGreaterThan adds a ports filter condition.
func WithPortsGreaterThan(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpGreaterThan, strconv.Itoa(ports))
	}
}

// WithPortsGreaterThanOrEqual adds a ports filter condition.
func WithPortsGreaterThanOrEqual(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpGreaterThanOrEqual, strconv.Itoa(ports))
	}
}

// WithPortsLessThan adds a ports filter condition.
func WithPortsLessThan(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpLessThan, strconv.Itoa(ports))
	}
}

// WithPortsLessThanOrEqual adds a ports filter condition.
func WithPortsLessThanOrEqual(ports int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("ports", OpLessThanOrEqual, strconv.Itoa(ports))
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

// WithType adds a type filter condition.
func WithType(type_ string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("type", OpEqual, type_)
	}
}

// WithTypeLike adds a type filter condition using LIKE operator.
func WithTypeLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("type", OpLike, pattern)
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

// WithResult adds a result filter condition.
func WithResult(result string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result", OpEqual, result)
	}
}

// WithResultLike adds a result filter condition using LIKE operator.
func WithResultLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result", OpLike, pattern)
	}
}

// WithTotal adds a total filter condition.
func WithTotal(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpEqual, strconv.Itoa(total))
	}
}

// WithTotalGreaterThan adds a total filter condition.
func WithTotalGreaterThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThan, strconv.Itoa(total))
	}
}

// WithTotalGreaterThanOrEqual adds a total filter condition.
func WithTotalGreaterThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThanOrEqual, strconv.Itoa(total))
	}
}

// WithTotalLessThan adds a total filter condition.
func WithTotalLessThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThan, strconv.Itoa(total))
	}
}

// WithTotalLessThanOrEqual adds a total filter condition.
func WithTotalLessThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThanOrEqual, strconv.Itoa(total))
	}
}

// WithFalsePositive adds a false_positive filter condition.
func WithFalsePositive(falsePositive int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("false_positive", OpEqual, strconv.Itoa(falsePositive))
	}
}

// WithFalsePositiveGreaterThan adds a false_positive filter condition.
func WithFalsePositiveGreaterThan(falsePositive int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("false_positive", OpGreaterThan, strconv.Itoa(falsePositive))
	}
}

// WithFalsePositiveGreaterThanOrEqual adds a false_positive filter condition.
func WithFalsePositiveGreaterThanOrEqual(falsePositive int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("false_positive", OpGreaterThanOrEqual, strconv.Itoa(falsePositive))
	}
}

// WithFalsePositiveLessThan adds a false_positive filter condition.
func WithFalsePositiveLessThan(falsePositive int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("false_positive", OpLessThan, strconv.Itoa(falsePositive))
	}
}

// WithFalsePositiveLessThanOrEqual adds a false_positive filter condition.
func WithFalsePositiveLessThanOrEqual(falsePositive int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("false_positive", OpLessThanOrEqual, strconv.Itoa(falsePositive))
	}
}

// WithLog adds a log filter condition.
func WithLog(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpEqual, strconv.Itoa(log))
	}
}

// WithLogGreaterThan adds a log filter condition.
func WithLogGreaterThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThan, strconv.Itoa(log))
	}
}

// WithLogGreaterThanOrEqual adds a log filter condition.
func WithLogGreaterThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThanOrEqual, strconv.Itoa(log))
	}
}

// WithLogLessThan adds a log filter condition.
func WithLogLessThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThan, strconv.Itoa(log))
	}
}

// WithLogLessThanOrEqual adds a log filter condition.
func WithLogLessThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThanOrEqual, strconv.Itoa(log))
	}
}

// WithLow adds a low filter condition.
func WithLow(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpEqual, strconv.Itoa(low))
	}
}

// WithLowGreaterThan adds a low filter condition.
func WithLowGreaterThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThan, strconv.Itoa(low))
	}
}

// WithLowGreaterThanOrEqual adds a low filter condition.
func WithLowGreaterThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThanOrEqual, strconv.Itoa(low))
	}
}

// WithLowLessThan adds a low filter condition.
func WithLowLessThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThan, strconv.Itoa(low))
	}
}

// WithLowLessThanOrEqual adds a low filter condition.
func WithLowLessThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThanOrEqual, strconv.Itoa(low))
	}
}

// WithMedium adds a medium filter condition.
func WithMedium(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpEqual, strconv.Itoa(medium))
	}
}

// WithMediumGreaterThan adds a medium filter condition.
func WithMediumGreaterThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThan, strconv.Itoa(medium))
	}
}

// WithMediumGreaterThanOrEqual adds a medium filter condition.
func WithMediumGreaterThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThanOrEqual, strconv.Itoa(medium))
	}
}

// WithMediumLessThan adds a medium filter condition.
func WithMediumLessThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThan, strconv.Itoa(medium))
	}
}

// WithMediumLessThanOrEqual adds a medium filter condition.
func WithMediumLessThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThanOrEqual, strconv.Itoa(medium))
	}
}

// WithHigh adds a high filter condition.
func WithHigh(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpEqual, strconv.Itoa(high))
	}
}

// WithHighGreaterThan adds a high filter condition.
func WithHighGreaterThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThan, strconv.Itoa(high))
	}
}

// WithHighGreaterThanOrEqual adds a high filter condition.
func WithHighGreaterThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThanOrEqual, strconv.Itoa(high))
	}
}

// WithHighLessThan adds a high filter condition.
func WithHighLessThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThan, strconv.Itoa(high))
	}
}

// WithHighLessThanOrEqual adds a high filter condition.
func WithHighLessThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThanOrEqual, strconv.Itoa(high))
	}
}

// WithCritical adds a critical filter condition.
func WithCritical(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpEqual, strconv.Itoa(critical))
	}
}

// WithCriticalGreaterThan adds a critical filter condition.
func WithCriticalGreaterThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThan, strconv.Itoa(critical))
	}
}

// WithCriticalGreaterThanOrEqual adds a critical filter condition.
func WithCriticalGreaterThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThanOrEqual, strconv.Itoa(critical))
	}
}

// WithCriticalLessThan adds a critical filter condition.
func WithCriticalLessThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThan, strconv.Itoa(critical))
	}
}

// WithCriticalLessThanOrEqual adds a critical filter condition.
func WithCriticalLessThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThanOrEqual, strconv.Itoa(critical))
	}
}

// WithResultHosts adds a result_hosts filter condition.
func WithResultHosts(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpEqual, strconv.Itoa(resultHosts))
	}
}

// WithResultHostsGreaterThan adds a result_hosts filter condition.
func WithResultHostsGreaterThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThan, strconv.Itoa(resultHosts))
	}
}

// WithResultHostsGreaterThanOrEqual adds a result_hosts filter condition.
func WithResultHostsGreaterThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThanOrEqual, strconv.Itoa(resultHosts))
	}
}

// WithResultHostsLessThan adds a result_hosts filter condition.
func WithResultHostsLessThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThan, strconv.Itoa(resultHosts))
	}
}

// WithResultHostsLessThanOrEqual adds a result_hosts filter condition.
func WithResultHostsLessThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThanOrEqual, strconv.Itoa(resultHosts))
	}
}

// WithSchedule adds a schedule filter condition.
func WithSchedule(schedule string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("schedule", OpEqual, schedule)
	}
}

// WithScheduleLike adds a schedule filter condition using LIKE operator.
func WithScheduleLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("schedule", OpLike, pattern)
	}
}

// WithNextDueAfter adds a next_due filter condition.
func WithNextDueAfter(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("next_due", OpGreaterThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithNextDueBefore adds a next_due filter condition.
func WithNextDueBefore(t time.Time) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("next_due", OpLessThan, t.Format("2006-01-02T15:04:05Z"))
	}
}

// WithNextDueBetween adds a next_due filter condition.
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
