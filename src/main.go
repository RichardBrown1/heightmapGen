package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
		scanner.Split(bufio.ScanWords)

		//Get ESRIInfo

		//Getting right hand side of col - should validate but esrigrids are standardised somewhat
		scanner.Scan()
		scanner.Scan()
		map1.ncols, err = strconv.Atoi(scanner.Text())

		scanner.Scan()
		scanner.Scan()
		map1.nrows, err = strconv.Atoi(scanner.Text())

		//TODO: There are xllcenter and yllcenter in some esri grids
		scanner.Scan()
		scanner.Scan()
		map1.xllcorner, err = strconv.ParseFloat(scanner.Text(), 64)

		scanner.Scan()
		scanner.Scan()
		map1.yllcorner, err = strconv.ParseFloat(scanner.Text(), 64)

		scanner.Scan()
		scanner.Scan()
		map1.cellsize, err = ParseFloat32(scanner.Text())

		// //nodata_value can be missing depending on implementation
		if scanner.Text() == "nodata_value" {
			scanner.Scan()
			map1.noDataValue, err = ParseFloat32(scanner.Text())
		} else {
			map1.noDataValue = noDataValueDefault
		}

		fmt.Println(map1.ncols, map1.nrows, map1.xllcorner, map1.noDataValue)

	}
}
