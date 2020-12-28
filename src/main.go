package main

import (
	"fmt"
	"os"
)

// type esriGrid struct {
// 	ncols       int
// 	nrows       int
// 	xllcorner   float64
// 	yllcorner   float64
// 	cellsize    float32
// 	noDataValue float32
// 	grid        [][]float32
// }

func main() {

	root := os.Args[1]

	fmt.Println("starting program at root path: " + root)

	//find ascii files
	allASCIIFiles := GetAllASCIIFiles(root)

	//fmt.Println(allASCIIFiles)
	esriGrids := GenerateEsriGrids(allASCIIFiles)
	_ = esriGrids
}
