package filtering

// WithReportID adds a report_id filter condition.
func WithReportID(reportID string) Option {
	return WithField("report_id", reportID)
}

// WithReportName adds a report name filter condition.
func WithReportName(name string) Option {
	return WithField("report_name", name)
}

// WithReportNameLike adds a report name filter condition using LIKE operator.
func WithReportNameLike(pattern string) Option {
	return WithFieldLike("report_name", pattern)
}

// WithReportFormat adds a report format filter condition.
func WithReportFormat(format string) Option {
	return WithField("report_format", format)
}

// WithReportFormatLike adds a report format filter condition using LIKE operator.
func WithReportFormatLike(pattern string) Option {
	return WithFieldLike("report_format", pattern)
}

// WithReportTask adds a task filter condition.
func WithReportTask(task string) Option {
	return WithField("task", task)
}

// WithReportTaskLike adds a task filter condition using LIKE operator.
func WithReportTaskLike(pattern string) Option {
	return WithFieldLike("task", pattern)
}

// WithReportScanner adds a scanner filter condition.
func WithReportScanner(scanner string) Option {
	return WithField("scanner", scanner)
}

// WithReportScannerLike adds a scanner filter condition using LIKE operator.
func WithReportScannerLike(pattern string) Option {
	return WithFieldLike("scanner", pattern)
}

// WithReportHosts adds a hosts filter condition.
func WithReportHosts(hosts int) Option {
	return WithNumericField("hosts", hosts)
}

// WithReportHostsGreaterThan adds a hosts filter condition.
func WithReportHostsGreaterThan(hosts int) Option {
	return WithNumericFieldGreaterThan("hosts", hosts)
}

// WithReportHostsGreaterThanOrEqual adds a hosts filter condition.
func WithReportHostsGreaterThanOrEqual(hosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("hosts", hosts)
}

// WithReportHostsLessThan adds a hosts filter condition.
func WithReportHostsLessThan(hosts int) Option {
	return WithNumericFieldLessThan("hosts", hosts)
}

// WithReportHostsLessThanOrEqual adds a hosts filter condition.
func WithReportHostsLessThanOrEqual(hosts int) Option {
	return WithNumericFieldLessThanOrEqual("hosts", hosts)
}

// WithReportResultHosts adds a result_hosts filter condition.
func WithReportResultHosts(resultHosts int) Option {
	return WithNumericField("result_hosts", resultHosts)
}

// WithReportResultHostsGreaterThan adds a result_hosts filter condition.
func WithReportResultHostsGreaterThan(resultHosts int) Option {
	return WithNumericFieldGreaterThan("result_hosts", resultHosts)
}

// WithReportResultHostsGreaterThanOrEqual adds a result_hosts filter condition.
func WithReportResultHostsGreaterThanOrEqual(resultHosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("result_hosts", resultHosts)
}

// WithReportResultHostsLessThan adds a result_hosts filter condition.
func WithReportResultHostsLessThan(resultHosts int) Option {
	return WithNumericFieldLessThan("result_hosts", resultHosts)
}

// WithReportResultHostsLessThanOrEqual adds a result_hosts filter condition.
func WithReportResultHostsLessThanOrEqual(resultHosts int) Option {
	return WithNumericFieldLessThanOrEqual("result_hosts", resultHosts)
}

// WithReportTotal adds a total filter condition.
func WithReportTotal(total int) Option {
	return WithNumericField("total", total)
}

// WithReportTotalGreaterThan adds a total filter condition.
func WithReportTotalGreaterThan(total int) Option {
	return WithNumericFieldGreaterThan("total", total)
}

// WithReportTotalGreaterThanOrEqual adds a total filter condition.
func WithReportTotalGreaterThanOrEqual(total int) Option {
	return WithNumericFieldGreaterThanOrEqual("total", total)
}

// WithReportTotalLessThan adds a total filter condition.
func WithReportTotalLessThan(total int) Option {
	return WithNumericFieldLessThan("total", total)
}

// WithReportTotalLessThanOrEqual adds a total filter condition.
func WithReportTotalLessThanOrEqual(total int) Option {
	return WithNumericFieldLessThanOrEqual("total", total)
}

// WithReportCritical adds a critical filter condition.
func WithReportCritical(critical int) Option {
	return WithNumericField("critical", critical)
}

// WithReportCriticalGreaterThan adds a critical filter condition.
func WithReportCriticalGreaterThan(critical int) Option {
	return WithNumericFieldGreaterThan("critical", critical)
}

// WithReportCriticalGreaterThanOrEqual adds a critical filter condition.
func WithReportCriticalGreaterThanOrEqual(critical int) Option {
	return WithNumericFieldGreaterThanOrEqual("critical", critical)
}

// WithReportCriticalLessThan adds a critical filter condition.
func WithReportCriticalLessThan(critical int) Option {
	return WithNumericFieldLessThan("critical", critical)
}

// WithReportCriticalLessThanOrEqual adds a critical filter condition.
func WithReportCriticalLessThanOrEqual(critical int) Option {
	return WithNumericFieldLessThanOrEqual("critical", critical)
}

// WithReportHigh adds a high filter condition.
func WithReportHigh(high int) Option {
	return WithNumericField("high", high)
}

// WithReportHighGreaterThan adds a high filter condition.
func WithReportHighGreaterThan(high int) Option {
	return WithNumericFieldGreaterThan("high", high)
}

// WithReportHighGreaterThanOrEqual adds a high filter condition.
func WithReportHighGreaterThanOrEqual(high int) Option {
	return WithNumericFieldGreaterThanOrEqual("high", high)
}

// WithReportHighLessThan adds a high filter condition.
func WithReportHighLessThan(high int) Option {
	return WithNumericFieldLessThan("high", high)
}

// WithReportHighLessThanOrEqual adds a high filter condition.
func WithReportHighLessThanOrEqual(high int) Option {
	return WithNumericFieldLessThanOrEqual("high", high)
}

// WithReportMedium adds a medium filter condition.
func WithReportMedium(medium int) Option {
	return WithNumericField("medium", medium)
}

// WithReportMediumGreaterThan adds a medium filter condition.
func WithReportMediumGreaterThan(medium int) Option {
	return WithNumericFieldGreaterThan("medium", medium)
}

// WithReportMediumGreaterThanOrEqual adds a medium filter condition.
func WithReportMediumGreaterThanOrEqual(medium int) Option {
	return WithNumericFieldGreaterThanOrEqual("medium", medium)
}

// WithReportMediumLessThan adds a medium filter condition.
func WithReportMediumLessThan(medium int) Option {
	return WithNumericFieldLessThan("medium", medium)
}

// WithReportMediumLessThanOrEqual adds a medium filter condition.
func WithReportMediumLessThanOrEqual(medium int) Option {
	return WithNumericFieldLessThanOrEqual("medium", medium)
}

// WithReportLow adds a low filter condition.
func WithReportLow(low int) Option {
	return WithNumericField("low", low)
}

// WithReportLowGreaterThan adds a low filter condition.
func WithReportLowGreaterThan(low int) Option {
	return WithNumericFieldGreaterThan("low", low)
}

// WithReportLowGreaterThanOrEqual adds a low filter condition.
func WithReportLowGreaterThanOrEqual(low int) Option {
	return WithNumericFieldGreaterThanOrEqual("low", low)
}

// WithReportLowLessThan adds a low filter condition.
func WithReportLowLessThan(low int) Option {
	return WithNumericFieldLessThan("low", low)
}

// WithReportLowLessThanOrEqual adds a low filter condition.
func WithReportLowLessThanOrEqual(low int) Option {
	return WithNumericFieldLessThanOrEqual("low", low)
}

// WithReportLog adds a log filter condition.
func WithReportLog(log int) Option {
	return WithNumericField("log", log)
}

// WithReportLogGreaterThan adds a log filter condition.
func WithReportLogGreaterThan(log int) Option {
	return WithNumericFieldGreaterThan("log", log)
}

// WithReportLogGreaterThanOrEqual adds a log filter condition.
func WithReportLogGreaterThanOrEqual(log int) Option {
	return WithNumericFieldGreaterThanOrEqual("log", log)
}

// WithReportLogLessThan adds a log filter condition.
func WithReportLogLessThan(log int) Option {
	return WithNumericFieldLessThan("log", log)
}

// WithReportLogLessThanOrEqual adds a log filter condition.
func WithReportLogLessThanOrEqual(log int) Option {
	return WithNumericFieldLessThanOrEqual("log", log)
}

// WithReportAlterable adds an alterable filter condition.
func WithReportAlterable(alterable bool) Option {
	return WithBooleanField("alterable", alterable)
}

// WithReportWritable adds a writable filter condition.
func WithReportWritable(writable bool) Option {
	return WithBooleanField("writable", writable)
}

// WithReportInUse adds an in_use filter condition.
func WithReportInUse(inUse bool) Option {
	return WithBooleanField("in_use", inUse)
}

// WithReportTrash adds a trash filter condition.
func WithReportTrash(trash bool) Option {
	return WithBooleanField("trash", trash)
}

// WithReport adds a report filter condition.
func WithReport(report string) Option {
	return WithField("report", report)
}

// WithReportLike adds a report filter condition using LIKE operator.
func WithReportLike(pattern string) Option {
	return WithFieldLike("report", pattern)
}

// WithTotal adds a total filter condition.
func WithTotal(total int) Option {
	return WithNumericField("total", total)
}

// WithTotalGreaterThan adds a total filter condition.
func WithTotalGreaterThan(total int) Option {
	return WithNumericFieldGreaterThan("total", total)
}

// WithTotalGreaterThanOrEqual adds a total filter condition.
func WithTotalGreaterThanOrEqual(total int) Option {
	return WithNumericFieldGreaterThanOrEqual("total", total)
}

// WithTotalLessThan adds a total filter condition.
func WithTotalLessThan(total int) Option {
	return WithNumericFieldLessThan("total", total)
}

// WithTotalLessThanOrEqual adds a total filter condition.
func WithTotalLessThanOrEqual(total int) Option {
	return WithNumericFieldLessThanOrEqual("total", total)
}

// WithLog adds a log filter condition.
func WithLog(log int) Option {
	return WithNumericField("log", log)
}

// WithLogGreaterThan adds a log filter condition.
func WithLogGreaterThan(log int) Option {
	return WithNumericFieldGreaterThan("log", log)
}

// WithLogGreaterThanOrEqual adds a log filter condition.
func WithLogGreaterThanOrEqual(log int) Option {
	return WithNumericFieldGreaterThanOrEqual("log", log)
}

// WithLogLessThan adds a log filter condition.
func WithLogLessThan(log int) Option {
	return WithNumericFieldLessThan("log", log)
}

// WithLogLessThanOrEqual adds a log filter condition.
func WithLogLessThanOrEqual(log int) Option {
	return WithNumericFieldLessThanOrEqual("log", log)
}

// WithLow adds a low filter condition.
func WithLow(low int) Option {
	return WithNumericField("low", low)
}

// WithLowGreaterThan adds a low filter condition.
func WithLowGreaterThan(low int) Option {
	return WithNumericFieldGreaterThan("low", low)
}

// WithLowGreaterThanOrEqual adds a low filter condition.
func WithLowGreaterThanOrEqual(low int) Option {
	return WithNumericFieldGreaterThanOrEqual("low", low)
}

// WithLowLessThan adds a low filter condition.
func WithLowLessThan(low int) Option {
	return WithNumericFieldLessThan("low", low)
}

// WithLowLessThanOrEqual adds a low filter condition.
func WithLowLessThanOrEqual(low int) Option {
	return WithNumericFieldLessThanOrEqual("low", low)
}

// WithMedium adds a medium filter condition.
func WithMedium(medium int) Option {
	return WithNumericField("medium", medium)
}

// WithMediumGreaterThan adds a medium filter condition.
func WithMediumGreaterThan(medium int) Option {
	return WithNumericFieldGreaterThan("medium", medium)
}

// WithMediumGreaterThanOrEqual adds a medium filter condition.
func WithMediumGreaterThanOrEqual(medium int) Option {
	return WithNumericFieldGreaterThanOrEqual("medium", medium)
}

// WithMediumLessThan adds a medium filter condition.
func WithMediumLessThan(medium int) Option {
	return WithNumericFieldLessThan("medium", medium)
}

// WithMediumLessThanOrEqual adds a medium filter condition.
func WithMediumLessThanOrEqual(medium int) Option {
	return WithNumericFieldLessThanOrEqual("medium", medium)
}

// WithHigh adds a high filter condition.
func WithHigh(high int) Option {
	return WithNumericField("high", high)
}

// WithHighGreaterThan adds a high filter condition.
func WithHighGreaterThan(high int) Option {
	return WithNumericFieldGreaterThan("high", high)
}

// WithHighGreaterThanOrEqual adds a high filter condition.
func WithHighGreaterThanOrEqual(high int) Option {
	return WithNumericFieldGreaterThanOrEqual("high", high)
}

// WithHighLessThan adds a high filter condition.
func WithHighLessThan(high int) Option {
	return WithNumericFieldLessThan("high", high)
}

// WithHighLessThanOrEqual adds a high filter condition.
func WithHighLessThanOrEqual(high int) Option {
	return WithNumericFieldLessThanOrEqual("high", high)
}

// WithCritical adds a critical filter condition.
func WithCritical(critical int) Option {
	return WithNumericField("critical", critical)
}

// WithCriticalGreaterThan adds a critical filter condition.
func WithCriticalGreaterThan(critical int) Option {
	return WithNumericFieldGreaterThan("critical", critical)
}

// WithCriticalGreaterThanOrEqual adds a critical filter condition.
func WithCriticalGreaterThanOrEqual(critical int) Option {
	return WithNumericFieldGreaterThanOrEqual("critical", critical)
}

// WithCriticalLessThan adds a critical filter condition.
func WithCriticalLessThan(critical int) Option {
	return WithNumericFieldLessThan("critical", critical)
}

// WithCriticalLessThanOrEqual adds a critical filter condition.
func WithCriticalLessThanOrEqual(critical int) Option {
	return WithNumericFieldLessThanOrEqual("critical", critical)
}

// WithResultHosts adds a result_hosts filter condition.
func WithResultHosts(resultHosts int) Option {
	return WithNumericField("result_hosts", resultHosts)
}

// WithResultHostsGreaterThan adds a result_hosts filter condition.
func WithResultHostsGreaterThan(resultHosts int) Option {
	return WithNumericFieldGreaterThan("result_hosts", resultHosts)
}

// WithResultHostsGreaterThanOrEqual adds a result_hosts filter condition.
func WithResultHostsGreaterThanOrEqual(resultHosts int) Option {
	return WithNumericFieldGreaterThanOrEqual("result_hosts", resultHosts)
}

// WithResultHostsLessThan adds a result_hosts filter condition.
func WithResultHostsLessThan(resultHosts int) Option {
	return WithNumericFieldLessThan("result_hosts", resultHosts)
}

// WithResultHostsLessThanOrEqual adds a result_hosts filter condition.
func WithResultHostsLessThanOrEqual(resultHosts int) Option {
	return WithNumericFieldLessThanOrEqual("result_hosts", resultHosts)
}

// WithSchedule adds a schedule filter condition.
func WithSchedule(schedule string) Option {
	return WithField("schedule", schedule)
}

// WithScheduleLike adds a schedule filter condition using LIKE operator.
func WithScheduleLike(pattern string) Option {
	return WithFieldLike("schedule", pattern)
}
