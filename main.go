package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	path := os.Args[1]

	fmt.Println("starting program at path: " + path)
}
