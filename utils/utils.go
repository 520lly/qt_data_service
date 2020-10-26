package utils

import(
   "os"
   "log"
)

func EnsurePathExist(path string)(bool, error) {
   _, err := os.Stat(path)
   if err == nil {
      return true, nil
   }
   if os.IsNotExist(err) {
      err := os.MkdirAll(path, 0755)
      if err == nil {
         return true, nil
      } else {
         log.Println(err)
         return false, err
      }
   }
   return false, err
}
