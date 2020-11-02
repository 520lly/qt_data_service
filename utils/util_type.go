package utils

import(
   "fmt"
)

func ConvertInterface2String(i interface{}) string {
   return fmt.Sprintf("%s", i)
}
