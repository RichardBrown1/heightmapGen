package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const noDataValueDefault = -9999.0

//EsriGrid that holds the ASCII data
type EsriGrid struct {
	ncols       int
	nrows       int
	xllcorner   float64
	yllcorner   float64
	cellsize    float32
	noDataValue float32
	grid        [][]float32
}

//GenerateEsriGrids of source files
func GenerateEsriGrids(ASCIIFilePaths []string) (esriGrids []EsriGrid) {
	//get files

	for _, fileName := range ASCIIFilePaths {
		fmt.Println(fileName)

		file, err := os.Open(fileName)
		defer file.Close()
		Check(err)

		var map1 EsriGrid

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanWords)

		getEsriInfo(&map1, scanner)
		getEsriGrid(&map1, scanner)

		esriGrids = append(esriGrids, map1)
	}
	return esriGrids
}

func getEsriGrid(eg *EsriGrid, s *bufio.Scanner) {
	eg.grid = make([][]float32, eg.nrows)
	for r := 1; r < eg.nrows; r++ {
		eg.grid[r] = make([]float32, eg.ncols)
		for c := 1; c < eg.ncols; c++ {
			eg.grid[r][c] = ParseFloat32(s.Text())
			if !s.Scan() {
				break
			}
		}
	}
}

func getEsriInfo(eg *EsriGrid, s *bufio.Scanner) {
	var err error

	//Getting right hand side of col - should validate but EsriGrids are standardised somewhat
	skipAndScan(s, 1)
	eg.ncols, err = strconv.Atoi(s.Text())
	Check(err)

	skipAndScan(s, 1)
	eg.nrows, err = strconv.Atoi(s.Text())
	Check(err)

	//TODO: There are xllcenter and yllcenter in some esri grids
	skipAndScan(s, 1)
	eg.xllcorner, err = strconv.ParseFloat(s.Text(), 64)
	Check(err)

	skipAndScan(s, 1)
	eg.yllcorner, err = strconv.ParseFloat(s.Text(), 64)
	Check(err)

	skipAndScan(s, 1)
	eg.cellsize = ParseFloat32(s.Text())

	//nodata_value can be missing depending on implementation -- wikipedia said so
	s.Scan()
	if s.Text() == "nodata_value" {
		s.Scan()
		eg.noDataValue = ParseFloat32(s.Text())
	} else {
		fmt.Println("'", s.Text(), "'")
		eg.noDataValue = noDataValueDefault
	}

	//fmt.Println(" cols", eg.ncols,
	// "\n rows", eg.nrows,
	// "\n xllcorner", eg.xllcorner,
	// "\n yllcorner", eg.yllcorner,
	// "\n cellSize", eg.cellsize,
	// "\n noDataVal", eg.noDataValue)

}

//GetAllASCIIFiles returns files, err
func GetAllASCIIFiles(path string) []string {
	var files []string
	parentFolder := path
	fileInfo, err := ioutil.ReadDir(path)
	Check(err)

	if runtime.GOOS == "windows" {
		parentFolder += "\\"
	} else {
		parentFolder += "/"
	}

	for _, file := range fileInfo {

		if strings.HasSuffix(file.Name(), ".asc") {
			files = append(files, parentFolder+file.Name())
		}

		if file.IsDir() {
			var subdirFiles []string
			subdirFiles = GetAllASCIIFiles(parentFolder + file.Name())
			files = append(files, subdirFiles...)
		}
	}
	return files
}

//skipAndScan (skips n scans on s)
func skipAndScan(s *bufio.Scanner, n int) {
	for n >= 0 {
		s.Scan()
		n--
	}
}
