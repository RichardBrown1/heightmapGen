package main

import (
	"bufio"
	"fmt"
	"os"
)

type esriGrid struct {
	ncols       int
	nrows       int
	xllcorner   float64
	yllcorner   float64
	cellsize    float32
	NodataValue float32
}

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

		scanner := bufio.NewScanner(file)
		for i := 0; i < 6; i++ {
			scanner.Scan()
			fmt.Println(scanner.Text())

		}
	}
}
