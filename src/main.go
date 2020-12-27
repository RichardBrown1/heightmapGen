package main

import (
	"fmt"
	"os"
)

func main() {
	root := os.Args[1]

	fmt.Println("starting program at root path: " + root)

	//find ascii files
	allASCIIFiles, err := GetAllASCIIFiles(root)
	Check(err)
	//fmt.Println(allASCIIFiles)

	//get files
	for _, fileName := range allASCIIFiles {
		fmt.Println(fileName)

		file, err := os.Open(fileName)
		Check(err)
		defer file.Close()

	}
}
