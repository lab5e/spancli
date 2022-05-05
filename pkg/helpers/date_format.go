package helpers

import (
	"strconv"
	"time"
)

// DateFormat formats a date to RFC3339 if the date formatting flag is set
func DateFormat(dateStr string, numeric bool) string {
	if numeric {
		return dateStr
	}
	if dateStr == "" {
		return "-"
	}
	val, err := strconv.ParseInt(dateStr, 10, 64)
	if err == nil {
		return time.UnixMilli(val).Format(time.RFC3339)
	}

	return "(error)"
}
