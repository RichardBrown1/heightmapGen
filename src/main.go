package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type esriGrid struct {
	ncols       int
	nrows       int
	xllcorner   float64
	yllcorner   float64
	cellsize    float32
	noDataValue float32
	grid        [][]float32
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

		var map1 esriGrid

		scanner := bufio.NewScanner(file)
		for i := 0; i < 6; i++ {
			scanner.Scan()

			wordSep := strings.Fields(scanner.Text())[1]

			switch i {
			case 0:
				map1.ncols, err = strconv.Atoi(wordSep)
			case 1:
				map1.nrows, err = strconv.Atoi(wordSep)
			case 2:
				map1.xllcorner, err = strconv.ParseFloat(wordSep, 64)
			case 3:
				map1.yllcorner, err = strconv.ParseFloat(wordSep, 64)
			case 4:
				map1.cellsize, err = ParseFloat32(wordSep)
			case 5: //This can be missing depending on implementation
				map1.noDataValue, err = ParseFloat32(wordSep)
			case default:

			}
		}

		fmt.Println(map1.ncols, map1.nrows, map1.xllcorner)
	}
}
