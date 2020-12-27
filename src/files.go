package main

import (
	"io/ioutil"
	"runtime"
	"strings"
)

//GetAllAsciiFiles returns files, err
func GetAllAsciiFiles(path string) ([]string, error) {
	var files []string
	parentFolder := path
	fileInfo, err := ioutil.ReadDir(path)

	if runtime.GOOS == "windows" {
		parentFolder += "\\"
	} else {
		parentFolder += "/"
	}

	for _, file := range fileInfo {

		if strings.HasSuffix(file.Name(), ".asc") {
			files = append(files, parentFolder+file.Name())
		}

		if file.IsDir() {
			var subdirFiles []string
			subdirFiles, err = GetAllAsciiFiles(parentFolder + file.Name())
			files = append(files, subdirFiles...)
		}
	}
	return files, err
}
