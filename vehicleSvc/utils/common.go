package utils

import "time"

func GetTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
