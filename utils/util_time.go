package utils

import (
   "time"
   //"log"
)

const (
	format string = "20060102"
)

func GetTodayString(format string) (time.Time, string) {
	today := time.Now()
	ts := today.Format(format)
	return today, ts
}

func CalcDateDiffByDay(start string, end string) int{
	s, _ := time.Parse(format, start)
	e, _ := time.Parse(format, end)
	diff := e.Sub(s)
	return int(diff.Hours() / 24)
}

func AddDays2Date(start string, y int, m int, d int) string {
	s, _ := time.Parse(format, start)
   //duration := time.Duration(diff)
   //e := s.Add(time.Hour * 24 * duration)
   e := s.AddDate(y, m, d)
   end := e.Format(format)
   return end
}

//func SubDays2Date(start string, diff int64) string {
	//s, _ := time.Parse(format, start)
   //duration := time.Duration(diff)
   //e := s.Sub(time.Hour * 24 * duration)
   //end := e.Format(format)
   //return end
//}


func IsDateAfter(end1 string, end2 string) bool {
	e1, _ := time.Parse(format, end1)
	e2, _ := time.Parse(format, end2)
	diff := e1.Sub(e2)
   return (diff > 0)
}

func IsDateBefore(d1 string, d2 string) bool {
	date1, _ := time.Parse(format, d1)
	date2, _ := time.Parse(format, d2)
	diff := date1.Sub(date2)
   return (diff < 0)
}


