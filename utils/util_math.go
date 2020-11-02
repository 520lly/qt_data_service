package utils

import (
   "math"
)

const (
   MAX_ITEMS_DAILY int = 5000
)

func CalcRoundsCeil(diff int) int{
   rounds := math.Ceil(float64(diff)/float64(MAX_ITEMS_DAILY))
   return int(rounds)
}

func Reverse2DArray(s [][]interface{}) {
   for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
      s[i], s[j] = s[j], s[i]
   }
}
