package utils

import (
   "time"
   "log"
)

const (
	DateFormat1 string = "20060102"
	DateFormat2 string = "2006-01-02"
)

func GetTodayString(format string) (time.Time, string) {
	today := time.Now()
	ts := today.Format(format)
	return today, ts
}

func CalcDateDiffByDay(format string, start string, end string) int{
	s, _ := time.Parse(format, start)
	e, _ := time.Parse(format, end)
	diff := e.Sub(s)
	return int(diff.Hours() / 24)
}

func AddDays2Date(format string, start string, y int, m int, d int) string {
	s, _ := time.Parse(format, start)
   //duration := time.Duration(diff)
   //e := s.Add(time.Hour * 24 * duration)
   e := s.AddDate(y, m, d)
   end := e.Format(format)
   log.Printf("start:%s format:%s s:%s end:%s d:%d", start, format, s, end, d)
   return end
}

func IsDateAfter(format string, end1 string, end2 string) bool {
	e1, _ := time.Parse(format, end1)
	e2, _ := time.Parse(format, end2)
	diff := e1.Sub(e2)
   return (diff > 0)
}

func IsDateAfterToday(format string, d string) bool {
	today := time.Now()
	then, _ := time.Parse(format, d)
	diff := today.Sub(then)
   return (diff > 0)
}

func IsDateBefore(format string, d1 string, d2 string) bool {
	date1, _ := time.Parse(format, d1)
	date2, _ := time.Parse(format, d2)
	diff := date1.Sub(date2)
   return (diff < 0)
}


