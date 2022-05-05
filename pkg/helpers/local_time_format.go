package helpers

import (
	"strconv"
	"time"
)

const timeFmt = "2006-01-02 15:04:05"

func msSinceEpochToTime(ts string) (int64, time.Time) {
	r, err := strconv.ParseInt(ts, 10, 63)
	if err != nil {
		return time.Now().UnixNano() / int64(time.Millisecond), time.Now()
	}
	return r, time.Unix(0, r*int64(time.Millisecond))
}

// LocalTimeFormat formats a timestamp string into local time
func LocalTimeFormat(ts string) string {
	_, t := msSinceEpochToTime(ts)
	return t.Local().Format(timeFmt)
}
