package helpers

// EllipsisString returns a string no longer than n. If the string is truncated an ellipsis is added
func EllipsisString(s string, n int) string {
	if len(s) > n && len(s) > 3 {
		return s[:n-3] + "..."
	}
	return s
}
