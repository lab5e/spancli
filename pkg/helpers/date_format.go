package helpers

import (
	"strconv"
	"time"
)

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
