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
		SkipAndScan(scanner, 1)
		map1.ncols, err = strconv.Atoi(scanner.Text())

		SkipAndScan(scanner, 1)
		map1.nrows, err = strconv.Atoi(scanner.Text())

		//map1.grid = map1.grid[:map1.nrows][:map1.ncols]

		//TODO: There are xllcenter and yllcenter in some esri grids
		SkipAndScan(scanner, 1)
		map1.xllcorner, err = strconv.ParseFloat(scanner.Text(), 64)

		SkipAndScan(scanner, 1)
		map1.yllcorner, err = strconv.ParseFloat(scanner.Text(), 64)

		SkipAndScan(scanner, 1)
		map1.cellsize = ParseFloat32(scanner.Text())

		// //nodata_value can be missing depending on implementation
		scanner.Scan()
		if scanner.Text() == "nodata_value" {
			scanner.Scan()
			map1.noDataValue = ParseFloat32(scanner.Text())
		} else {
			fmt.Println("'", scanner.Text(), "'")
			map1.noDataValue = noDataValueDefault
		}

		map1.grid = make([][]float32, map1.nrows)

		for r := 1; r < map1.nrows; r++ {
			map1.grid[r] = make([]float32, map1.ncols)
			for c := 1; c < map1.ncols; c++ {
				map1.grid[r][c] = ParseFloat32(scanner.Text())
				if !scanner.Scan() {
					break
				}
			}
		}
		fmt.Println(map1.ncols, map1.nrows, map1.xllcorner, map1.noDataValue)
	}
}
