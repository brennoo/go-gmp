package filtering

import (
	"fmt"
)

// WithReportID adds a report_id filter condition.
func WithReportID(reportID string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report_id", OpEqual, reportID)
	}
}

// WithReportName adds a report name filter condition.
func WithReportName(name string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report_name", OpEqual, name)
	}
}

// WithReportNameLike adds a report name filter condition using LIKE operator.
func WithReportNameLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report_name", OpLike, pattern)
	}
}

// WithReportFormat adds a report format filter condition.
func WithReportFormat(format string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report_format", OpEqual, format)
	}
}

// WithReportFormatLike adds a report format filter condition using LIKE operator.
func WithReportFormatLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("report_format", OpLike, pattern)
	}
}

// WithReportTask adds a task filter condition.
func WithReportTask(task string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpEqual, task)
	}
}

// WithReportTaskLike adds a task filter condition using LIKE operator.
func WithReportTaskLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("task", OpLike, pattern)
	}
}

// WithReportScanner adds a scanner filter condition.
func WithReportScanner(scanner string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpEqual, scanner)
	}
}

// WithReportScannerLike adds a scanner filter condition using LIKE operator.
func WithReportScannerLike(pattern string) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("scanner", OpLike, pattern)
	}
}

// WithReportHosts adds a hosts filter condition.
func WithReportHosts(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithReportHostsGreaterThan adds a hosts filter condition.
func WithReportHostsGreaterThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThan, fmt.Sprintf("%d", hosts))
	}
}

// WithReportHostsGreaterThanOrEqual adds a hosts filter condition.
func WithReportHostsGreaterThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithReportHostsLessThan adds a hosts filter condition.
func WithReportHostsLessThan(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThan, fmt.Sprintf("%d", hosts))
	}
}

// WithReportHostsLessThanOrEqual adds a hosts filter condition.
func WithReportHostsLessThanOrEqual(hosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("hosts", OpLessThanOrEqual, fmt.Sprintf("%d", hosts))
	}
}

// WithReportResultHosts adds a result_hosts filter condition.
func WithReportResultHosts(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpEqual, fmt.Sprintf("%d", resultHosts))
	}
}

// WithReportResultHostsGreaterThan adds a result_hosts filter condition.
func WithReportResultHostsGreaterThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThan, fmt.Sprintf("%d", resultHosts))
	}
}

// WithReportResultHostsGreaterThanOrEqual adds a result_hosts filter condition.
func WithReportResultHostsGreaterThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", resultHosts))
	}
}

// WithReportResultHostsLessThan adds a result_hosts filter condition.
func WithReportResultHostsLessThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThan, fmt.Sprintf("%d", resultHosts))
	}
}

// WithReportResultHostsLessThanOrEqual adds a result_hosts filter condition.
func WithReportResultHostsLessThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThanOrEqual, fmt.Sprintf("%d", resultHosts))
	}
}

// WithReportTotal adds a total filter condition.
func WithReportTotal(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpEqual, fmt.Sprintf("%d", total))
	}
}

// WithReportTotalGreaterThan adds a total filter condition.
func WithReportTotalGreaterThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThan, fmt.Sprintf("%d", total))
	}
}

// WithReportTotalGreaterThanOrEqual adds a total filter condition.
func WithReportTotalGreaterThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThanOrEqual, fmt.Sprintf("%d", total))
	}
}

// WithReportTotalLessThan adds a total filter condition.
func WithReportTotalLessThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThan, fmt.Sprintf("%d", total))
	}
}

// WithReportTotalLessThanOrEqual adds a total filter condition.
func WithReportTotalLessThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThanOrEqual, fmt.Sprintf("%d", total))
	}
}

// WithReportCritical adds a critical filter condition.
func WithReportCritical(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpEqual, fmt.Sprintf("%d", critical))
	}
}

// WithReportCriticalGreaterThan adds a critical filter condition.
func WithReportCriticalGreaterThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThan, fmt.Sprintf("%d", critical))
	}
}

// WithReportCriticalGreaterThanOrEqual adds a critical filter condition.
func WithReportCriticalGreaterThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThanOrEqual, fmt.Sprintf("%d", critical))
	}
}

// WithReportCriticalLessThan adds a critical filter condition.
func WithReportCriticalLessThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThan, fmt.Sprintf("%d", critical))
	}
}

// WithReportCriticalLessThanOrEqual adds a critical filter condition.
func WithReportCriticalLessThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThanOrEqual, fmt.Sprintf("%d", critical))
	}
}

// WithReportHigh adds a high filter condition.
func WithReportHigh(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpEqual, fmt.Sprintf("%d", high))
	}
}

// WithReportHighGreaterThan adds a high filter condition.
func WithReportHighGreaterThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThan, fmt.Sprintf("%d", high))
	}
}

// WithReportHighGreaterThanOrEqual adds a high filter condition.
func WithReportHighGreaterThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThanOrEqual, fmt.Sprintf("%d", high))
	}
}

// WithReportHighLessThan adds a high filter condition.
func WithReportHighLessThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThan, fmt.Sprintf("%d", high))
	}
}

// WithReportHighLessThanOrEqual adds a high filter condition.
func WithReportHighLessThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThanOrEqual, fmt.Sprintf("%d", high))
	}
}

// WithReportMedium adds a medium filter condition.
func WithReportMedium(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpEqual, fmt.Sprintf("%d", medium))
	}
}

// WithReportMediumGreaterThan adds a medium filter condition.
func WithReportMediumGreaterThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThan, fmt.Sprintf("%d", medium))
	}
}

// WithReportMediumGreaterThanOrEqual adds a medium filter condition.
func WithReportMediumGreaterThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThanOrEqual, fmt.Sprintf("%d", medium))
	}
}

// WithReportMediumLessThan adds a medium filter condition.
func WithReportMediumLessThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThan, fmt.Sprintf("%d", medium))
	}
}

// WithReportMediumLessThanOrEqual adds a medium filter condition.
func WithReportMediumLessThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThanOrEqual, fmt.Sprintf("%d", medium))
	}
}

// WithReportLow adds a low filter condition.
func WithReportLow(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpEqual, fmt.Sprintf("%d", low))
	}
}

// WithReportLowGreaterThan adds a low filter condition.
func WithReportLowGreaterThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThan, fmt.Sprintf("%d", low))
	}
}

// WithReportLowGreaterThanOrEqual adds a low filter condition.
func WithReportLowGreaterThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThanOrEqual, fmt.Sprintf("%d", low))
	}
}

// WithReportLowLessThan adds a low filter condition.
func WithReportLowLessThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThan, fmt.Sprintf("%d", low))
	}
}

// WithReportLowLessThanOrEqual adds a low filter condition.
func WithReportLowLessThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThanOrEqual, fmt.Sprintf("%d", low))
	}
}

// WithReportLog adds a log filter condition.
func WithReportLog(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpEqual, fmt.Sprintf("%d", log))
	}
}

// WithReportLogGreaterThan adds a log filter condition.
func WithReportLogGreaterThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThan, fmt.Sprintf("%d", log))
	}
}

// WithReportLogGreaterThanOrEqual adds a log filter condition.
func WithReportLogGreaterThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThanOrEqual, fmt.Sprintf("%d", log))
	}
}

// WithReportLogLessThan adds a log filter condition.
func WithReportLogLessThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThan, fmt.Sprintf("%d", log))
	}
}

// WithReportLogLessThanOrEqual adds a log filter condition.
func WithReportLogLessThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThanOrEqual, fmt.Sprintf("%d", log))
	}
}

// WithReportAlterable adds an alterable filter condition.
func WithReportAlterable(alterable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("alterable", OpEqual, fmt.Sprintf("%t", alterable))
	}
}

// WithReportWritable adds a writable filter condition.
func WithReportWritable(writable bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("writable", OpEqual, fmt.Sprintf("%t", writable))
	}
}

// WithReportInUse adds an in_use filter condition.
func WithReportInUse(inUse bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("in_use", OpEqual, fmt.Sprintf("%t", inUse))
	}
}

// WithReportTrash adds a trash filter condition.
func WithReportTrash(trash bool) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("trash", OpEqual, fmt.Sprintf("%t", trash))
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

// WithTotal adds a total filter condition.
func WithTotal(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpEqual, fmt.Sprintf("%d", total))
	}
}

// WithTotalGreaterThan adds a total filter condition.
func WithTotalGreaterThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThan, fmt.Sprintf("%d", total))
	}
}

// WithTotalGreaterThanOrEqual adds a total filter condition.
func WithTotalGreaterThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpGreaterThanOrEqual, fmt.Sprintf("%d", total))
	}
}

// WithTotalLessThan adds a total filter condition.
func WithTotalLessThan(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThan, fmt.Sprintf("%d", total))
	}
}

// WithTotalLessThanOrEqual adds a total filter condition.
func WithTotalLessThanOrEqual(total int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("total", OpLessThanOrEqual, fmt.Sprintf("%d", total))
	}
}

// WithLog adds a log filter condition.
func WithLog(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpEqual, fmt.Sprintf("%d", log))
	}
}

// WithLogGreaterThan adds a log filter condition.
func WithLogGreaterThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThan, fmt.Sprintf("%d", log))
	}
}

// WithLogGreaterThanOrEqual adds a log filter condition.
func WithLogGreaterThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpGreaterThanOrEqual, fmt.Sprintf("%d", log))
	}
}

// WithLogLessThan adds a log filter condition.
func WithLogLessThan(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThan, fmt.Sprintf("%d", log))
	}
}

// WithLogLessThanOrEqual adds a log filter condition.
func WithLogLessThanOrEqual(log int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("log", OpLessThanOrEqual, fmt.Sprintf("%d", log))
	}
}

// WithLow adds a low filter condition.
func WithLow(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpEqual, fmt.Sprintf("%d", low))
	}
}

// WithLowGreaterThan adds a low filter condition.
func WithLowGreaterThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThan, fmt.Sprintf("%d", low))
	}
}

// WithLowGreaterThanOrEqual adds a low filter condition.
func WithLowGreaterThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpGreaterThanOrEqual, fmt.Sprintf("%d", low))
	}
}

// WithLowLessThan adds a low filter condition.
func WithLowLessThan(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThan, fmt.Sprintf("%d", low))
	}
}

// WithLowLessThanOrEqual adds a low filter condition.
func WithLowLessThanOrEqual(low int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("low", OpLessThanOrEqual, fmt.Sprintf("%d", low))
	}
}

// WithMedium adds a medium filter condition.
func WithMedium(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpEqual, fmt.Sprintf("%d", medium))
	}
}

// WithMediumGreaterThan adds a medium filter condition.
func WithMediumGreaterThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThan, fmt.Sprintf("%d", medium))
	}
}

// WithMediumGreaterThanOrEqual adds a medium filter condition.
func WithMediumGreaterThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpGreaterThanOrEqual, fmt.Sprintf("%d", medium))
	}
}

// WithMediumLessThan adds a medium filter condition.
func WithMediumLessThan(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThan, fmt.Sprintf("%d", medium))
	}
}

// WithMediumLessThanOrEqual adds a medium filter condition.
func WithMediumLessThanOrEqual(medium int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("medium", OpLessThanOrEqual, fmt.Sprintf("%d", medium))
	}
}

// WithHigh adds a high filter condition.
func WithHigh(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpEqual, fmt.Sprintf("%d", high))
	}
}

// WithHighGreaterThan adds a high filter condition.
func WithHighGreaterThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThan, fmt.Sprintf("%d", high))
	}
}

// WithHighGreaterThanOrEqual adds a high filter condition.
func WithHighGreaterThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpGreaterThanOrEqual, fmt.Sprintf("%d", high))
	}
}

// WithHighLessThan adds a high filter condition.
func WithHighLessThan(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThan, fmt.Sprintf("%d", high))
	}
}

// WithHighLessThanOrEqual adds a high filter condition.
func WithHighLessThanOrEqual(high int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("high", OpLessThanOrEqual, fmt.Sprintf("%d", high))
	}
}

// WithCritical adds a critical filter condition.
func WithCritical(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpEqual, fmt.Sprintf("%d", critical))
	}
}

// WithCriticalGreaterThan adds a critical filter condition.
func WithCriticalGreaterThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThan, fmt.Sprintf("%d", critical))
	}
}

// WithCriticalGreaterThanOrEqual adds a critical filter condition.
func WithCriticalGreaterThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpGreaterThanOrEqual, fmt.Sprintf("%d", critical))
	}
}

// WithCriticalLessThan adds a critical filter condition.
func WithCriticalLessThan(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThan, fmt.Sprintf("%d", critical))
	}
}

// WithCriticalLessThanOrEqual adds a critical filter condition.
func WithCriticalLessThanOrEqual(critical int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("critical", OpLessThanOrEqual, fmt.Sprintf("%d", critical))
	}
}

// WithResultHosts adds a result_hosts filter condition.
func WithResultHosts(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpEqual, fmt.Sprintf("%d", resultHosts))
	}
}

// WithResultHostsGreaterThan adds a result_hosts filter condition.
func WithResultHostsGreaterThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThan, fmt.Sprintf("%d", resultHosts))
	}
}

// WithResultHostsGreaterThanOrEqual adds a result_hosts filter condition.
func WithResultHostsGreaterThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpGreaterThanOrEqual, fmt.Sprintf("%d", resultHosts))
	}
}

// WithResultHostsLessThan adds a result_hosts filter condition.
func WithResultHostsLessThan(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThan, fmt.Sprintf("%d", resultHosts))
	}
}

// WithResultHostsLessThanOrEqual adds a result_hosts filter condition.
func WithResultHostsLessThanOrEqual(resultHosts int) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("result_hosts", OpLessThanOrEqual, fmt.Sprintf("%d", resultHosts))
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
