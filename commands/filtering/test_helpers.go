package filtering

// Helper function to check if a string contains a substring.
func containsHelper(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && containsHelper(s[1:], substr)
}
