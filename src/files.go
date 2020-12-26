package main

import (
	"io/ioutil"
)

// ListSubDirs returns file, err
func ListSubDirs(path string) ([]string, error) {
	fileInfo, err := ioutil.ReadDir(path)

	var files []string
	for _, file := range fileInfo {
		if file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, err
}
