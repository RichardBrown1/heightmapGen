package main

import (
	"io/ioutil"
	"strings"
)

//GetAllSubDirs returns files, err
func GetAllSubDirs(path string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(path)

	for _, file := range fileInfo {

		if strings.HasSuffix(file.Name(), ".asc") {
			files = append(files, path+"/"+file.Name())
		}

		if file.IsDir() {
			var subdirFiles []string
			subdirFiles, err = GetAllSubDirs(path + "/" + file.Name())
			files = append(files, subdirFiles...)
		}
	}
	return files, err
}
