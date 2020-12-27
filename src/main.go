package main

import (
	"fmt"
	"os"
)

func main() {
	root := os.Args[1]

	fmt.Println("starting program at root path: " + root)

	allASCIIFiles, err := GetAllAsciiFiles(root)
	Check(err)
	fmt.Println(allASCIIFiles)
}
