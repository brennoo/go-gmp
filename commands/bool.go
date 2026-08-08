package commands

// Bool returns a pointer to b, for setting a *bool request field to an
// explicit value (including false) rather than leaving it unset.
func Bool(b bool) *bool { return &b }
