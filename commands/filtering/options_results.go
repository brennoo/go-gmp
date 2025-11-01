package filtering

import (
	"fmt"
)

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
	return WithField("threat", threat)
}

// WithTrend adds a trend filter condition.
func WithTrend(trend string) Option {
	return WithField("trend", trend)
}

// WithQoD adds a quality of detection filter condition.
func WithQoD(qod float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpEqual, fmt.Sprintf("%.1f", qod))
	}
}

// WithQoDGreaterThan adds a quality of detection filter condition.
func WithQoDGreaterThan(qod float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpGreaterThan, fmt.Sprintf("%.1f", qod))
	}
}

// WithQoDGreaterThanOrEqual adds a quality of detection filter condition.
func WithQoDGreaterThanOrEqual(qod float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpGreaterThanOrEqual, fmt.Sprintf("%.1f", qod))
	}
}

// WithQoDLessThan adds a quality of detection filter condition.
func WithQoDLessThan(qod float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpLessThan, fmt.Sprintf("%.1f", qod))
	}
}

// WithQoDLessThanOrEqual adds a quality of detection filter condition.
func WithQoDLessThanOrEqual(qod float64) Option {
	return func(fb *FilterBuilder) {
		fb.AddCondition("qod", OpLessThanOrEqual, fmt.Sprintf("%.1f", qod))
	}
}

// WithHost adds a host filter condition.
func WithHost(host string) Option {
	return WithField("host", host)
}

// WithHostLike adds a host filter condition using LIKE operator.
func WithHostLike(pattern string) Option {
	return WithFieldLike("host", pattern)
}

// WithPort adds a port filter condition.
func WithPort(port string) Option {
	return WithField("port", port)
}

// WithPortLike adds a port filter condition using LIKE operator.
func WithPortLike(pattern string) Option {
	return WithFieldLike("port", pattern)
}

// WithNVT adds an NVT filter condition.
func WithNVT(nvt string) Option {
	return WithField("nvt", nvt)
}

// WithNVTLike adds an NVT filter condition using LIKE operator.
func WithNVTLike(pattern string) Option {
	return WithFieldLike("nvt", pattern)
}

// WithTask adds a task filter condition.
func WithTask(task string) Option {
	return WithField("task", task)
}

// WithTaskLike adds a task filter condition using LIKE operator.
func WithTaskLike(pattern string) Option {
	return WithFieldLike("task", pattern)
}

// WithResult adds a result filter condition.
func WithResult(result string) Option {
	return WithField("result", result)
}

// WithResultLike adds a result filter condition using LIKE operator.
func WithResultLike(pattern string) Option {
	return WithFieldLike("result", pattern)
}

// WithFalsePositive adds a false_positive filter condition.
func WithFalsePositive(fp bool) Option {
	return WithBooleanField("false_positive", fp)
}

// WithFalsePositiveGreaterThan adds a false_positive filter condition.
func WithFalsePositiveGreaterThan(count int) Option {
	return WithNumericFieldGreaterThan("false_positive", count)
}

// WithFalsePositiveGreaterThanOrEqual adds a false_positive filter condition.
func WithFalsePositiveGreaterThanOrEqual(count int) Option {
	return WithNumericFieldGreaterThanOrEqual("false_positive", count)
}

// WithFalsePositiveLessThan adds a false_positive filter condition.
func WithFalsePositiveLessThan(count int) Option {
	return WithNumericFieldLessThan("false_positive", count)
}

// WithFalsePositiveLessThanOrEqual adds a false_positive filter condition.
func WithFalsePositiveLessThanOrEqual(count int) Option {
	return WithNumericFieldLessThanOrEqual("false_positive", count)
}
