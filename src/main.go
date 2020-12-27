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
	const noDataValueDefault = -9999.0
	for _, fileName := range allASCIIFiles {
		fmt.Println(fileName)

		file, err := os.Open(fileName)
		Check(err)
		defer file.Close()

		var map1 esriGrid

		scanner := bufio.NewScanner(file)

		//Get ESRIInfo
		scanner.Scan()
		map1.ncols, err = strconv.Atoi(strings.Fields(scanner.Text())[1])

		scanner.Scan()
		map1.nrows, err = strconv.Atoi(strings.Fields(scanner.Text())[1])

		//TODO: There are xllcenter and yllcenter in some esri grids
		scanner.Scan()
		map1.xllcorner, err = strconv.ParseFloat(strings.Fields(scanner.Text())[1], 64)

		scanner.Scan()
		map1.yllcorner, err = strconv.ParseFloat(strings.Fields(scanner.Text())[1], 64)

		scanner.Scan()
		map1.cellsize, err = ParseFloat32(strings.Fields(scanner.Text())[1])

		//nodata_value can be missing depending on implementation
		scanner.Scan()
		if strings.Fields(scanner.Text())[0] == "nodata_value" {
			map1.noDataValue, err = ParseFloat32(strings.Fields(scanner.Text())[1])
		} else {
			map1.noDataValue = noDataValueDefault //Default
		}

		fmt.Println(map1.ncols, map1.nrows, map1.xllcorner, map1.noDataValue)
	}
}
