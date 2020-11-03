package utils

import (
   "regexp"
   "strings"
   "path/filepath"
)

func GetUpdateDateFromName(fileName string) string {
   re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
   submatchall := re.FindAllString(fileName, -1)
   if len(submatchall) >= 1 {
      return submatchall[0]
   }
   return ""
}

func IsFileSameFromFullPath(f string, fullpath string) bool {
   base := filepath.Base(fullpath)
   name := GetNameFromExtend(base)
   return f == name
}

func GetNameFromExtend(fullpath string) string {
   return strings.TrimSuffix(fullpath, filepath.Ext(fullpath))
}
