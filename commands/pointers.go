package commands

// Int returns a pointer to v, for setting an *int request field to an
// explicit value (including 0) rather than leaving it unset.
func Int(v int) *int {
	return &v
}

// Float64 returns a pointer to v, for setting a *float64 request field to an
// explicit value (including 0) rather than leaving it unset.
func Float64(v float64) *float64 {
	return &v
}
