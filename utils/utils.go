package utils

import (
	"log"
	"os"
	"path/filepath"
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
