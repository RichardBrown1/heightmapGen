package main

import (
	"io/ioutil"
	"strings"
)

//GetAllSubDirs returns files, err
func GetAllSubDirs(path string) ([]string, error) {
	fileInfo, err := ioutil.ReadDir(path)

	var files []string
	for _, file := range fileInfo {

		files = append(files, path+"/"+file.Name())

		if file.IsDir() || strings.HasSuffix(file.Name(), ".asc") {
			var subdirFiles []string

			if file.IsDir() {
				subdirFiles, err = GetAllSubDirs(path + "/" + file.Name())
			}
			files = append(files, subdirFiles...)
		}
	}
	return files, err
}

// // ListSubDirs returns file, err
// func ListSubDirs(path string) ([]string, error) {
// 	fileInfo, err := ioutil.ReadDir(path)

// 	var files []string
// 	for _, file := range fileInfo {
// 		if file.IsDir() {
// 			files = append(files, file.Name())
// 		}
// 	}
// 	return files, err
// }
