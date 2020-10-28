package utils

import (
	"log"
	"os"
	"path/filepath"
	"encoding/csv"
	"regexp"
)

func EnsurePathExist(path string) (bool, error) {
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

func CheckFileIsExist(fileName string) bool {
   var exist = true
   if _, err := os.Stat(fileName); os.IsNotExist(err) {
      exist = false
   }
   return exist
}

func FilteredSearchOfDirectoryTree(re *regexp.Regexp, dir string) (error, []string) {
   // Just a demo, this is how we capture the files that match the pattern.
	files := []string{}

	// Function variable that can be used to filter
	// files based on the pattern.
	// Note that it uses re internally to filter.
	// Also note that it populates the files variable with
	// the files that matches the pattern.
	walk := func(fn string, fi os.FileInfo, err error) error {
		if re.MatchString(fn) == false {
			return nil
		}
		if fi.IsDir() {
			log.Println(fn + string(os.PathSeparator))

		} else {
			log.Println(fn)
			files = append(files, fn)
		}
		return nil
	}
	filepath.Walk(dir, walk)
	log.Printf("Found %[1]d files.\n", len(files))
	return nil, files
}

func GetFileSize(file string) (int64, error) {
	fp, err := os.Stat(file)
	if err != nil {
		return -1, err
	}
   size := fp.Size()
	return size, nil
}

func OpenCsvFile(fileName string) (bool, *os.File) {
	var file *os.File
	var err1 error
	var isNew = false

	if CheckFileIsExist(fileName) {
      file, err1 = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	} else {
		file, err1 = os.Create(fileName)
		file.WriteString("\xEF\xBB\xBF") //for writing Chinese to csv file
		isNew = true
	}
	if err1 != nil {
		log.Println("unable to write file on filehook ", err1)
		panic(err1)
	}
	return isNew, file
}

func WriteData2CsvFile(csvw *csv.Writer, data []string) {
	if csvw == nil {
		return
	}
	//data := models.FieldSymbol
	csvw.Write(data)
	csvw.Flush()
}