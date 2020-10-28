package utils

import (
	"time"
)

func GetTodayString(format string) (time.Time, string){
   today := time.Now()
   ts := today.Format(format)
   return today,ts
}
