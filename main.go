package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	root := os.Args[1]

	fmt.Println("starting program at root path: " + root)

	fileInfo, err := ioutil.ReadDir(root)
	check(err)

	var files []string
	for _, file := range fileInfo {
		fmt.Println(file.Name())
		//files = append(files, file.Name())
	}

	fmt.Println(files)
}
