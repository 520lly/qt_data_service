package utils

import (
	"time"
)

func GetTodayString(format string) (time.Time, string) {
	today := time.Now()
	ts := today.Format(format)
	return today, ts
}

func CalcDateDiffByDay(start string, end string) int64 {
	format := "20060102"
	s, _ := time.Parse(format, start)
	e, _ := time.Parse(format, end)
	diff := e.Sub(s)
	return int64(diff.Hours() / 24)
}
