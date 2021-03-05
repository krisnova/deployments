package thenovadiary

import (
	"fmt"
	"time"
)

const (

	// TimeLayoutHuman is the time that is
	// relevant for me and you :)
	TimeLayoutHuman = "2006-01-02T15:04:05"

	// TimeLayout8601 is the ISO 8601 Standard
	TimeLayout8601 = "2006-01-02T15:04:05Z0700"

	// TimeLayoutTimeTime Is the format used by the Go
	// standard library time.Time String() method
	// Use this to decode a time.Time string -> *time.Time
	// instance with time.Parse()
	//
	// 2020-11-05 00:00:00 -0800 PST
	TimeLayoutTimeTime = "2006-01-02 15:04:05 -0700 PST"
)

// TimeDeltaDays will calculate the absolute
// value of the difference (delta) in days
// between two given points in time.
func TimeDeltaDays(t1, t2 time.Time) int {
	// T0 - - T1 - T2 - - Tn ->	true
	if t2.After(t1) {
		return int(t2.Sub(t1).Hours() / 24)
	}
	return int(t1.Sub(t2).Hours() / 24)
}

func TimeDeltaDaysFromNow(delta int) time.Time {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day()+delta, 0, 0, 0, 0, time.Local)
	return date
}

func TimeStringToTime(str string) (*time.Time, error) {
	time, err := time.Parse(TimeLayout8601, str)
	if err != nil {
		return nil, fmt.Errorf("time conversion error: %v", err)
	}
	return &time, nil
}

func TimeTimeToString(t time.Time) string {
	return t.String()
}

func TimeTomorrow() time.Time {
	return TimeDeltaDaysFromNow(1)
}

func TimeToday() time.Time {
	return TimeDeltaDaysFromNow(0)
}

func TimeYesterday() time.Time {
	return TimeDeltaDaysFromNow(-1)
}

func CronIsExpiredDays(timeToCheck time.Time, days int) bool {
	today := TimeToday()
	i := TimeDeltaDays(timeToCheck, today)
	if i > days {
		return true
	}
	return false
}
