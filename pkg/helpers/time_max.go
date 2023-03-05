package helpers

import "time"

var currentTime *time.Time

func TimeMax() time.Time {
	return time.Date(9998, time.Month(12), 31, 23, 59, 0, 0, time.Now().Location())
}

func Now() time.Time {
	if currentTime == nil {
		return time.Now()
	}
	return *currentTime
}

func SetTime(t time.Time) {
	currentTime = &t
}

func ResetTime() {
	currentTime = nil
}
