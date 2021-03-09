package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
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
	xllcorner   int
	yllcorner   int
	cellsize    int
	noDataValue float32
	grid        [][]float32
}

//GetEsriGridHighLow find highs and lows
func GetEsriGridHighLow(esriGrids []EsriGrid) (highY int, lowY int, highX int, lowX int, nRows int, nCols int, cellSize int) {
	//todo: stitch heightmaps together
	//for now assume cellsize is the same and its all from 1 data source
	lowX = math.MaxInt32
	lowY = math.MaxInt32
	highX = math.MaxInt32 * -1 //minFloat64
	highY = math.MaxInt32 * -1

	cellSize = esriGrids[0].cellsize
	nRows = esriGrids[0].nrows
	nCols = esriGrids[0].ncols
	for _, eg := range esriGrids {
		//find lowest xll and yll corner
		if eg.xllcorner < lowX {
			lowX = eg.xllcorner
		} else {
			if eg.xllcorner > highX {
				highX = eg.xllcorner
			}
		}
		if eg.yllcorner < lowY {
			lowY = eg.yllcorner
		} else {
			if eg.yllcorner > highY {
				highY = eg.yllcorner
			}
		}
		if eg.cellsize != cellSize || eg.nrows != nRows || eg.ncols != nCols {
			fmt.Println("cellsize, row or col counts are not consistent ... will terminate this program isnt good enough to handle that ")
			panic(1)
		}
	}

	fmt.Println(
		" lowX: ", lowX,
		"\n lowY: ", lowY,
		"\n highX: ", highX,
		"\n highY: ", highY,
		"\n nCols: ", nCols,
		"\n nRows: ", nRows,
		"\n Cellsize: ", cellSize)

	return highY, lowY, highX, lowX, nRows, nCols, cellSize
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

	fmt.Println(esriGrids[0].grid[0][0])

	return esriGrids
}

func getEsriGrid(eg *EsriGrid, s *bufio.Scanner) {
	eg.grid = make([][]float32, eg.nrows)
	for r := 0; r < eg.nrows; r++ {
		eg.grid[r] = make([]float32, eg.ncols)
		for c := 0; c < eg.ncols; c++ {
			eg.grid[r][c] = ParseFloat32(s.Text())
			//fmt.Println(eg.grid[r][c])
			if !s.Scan() {
				break
			}
		}
	}
}

func getEsriInfo(eg *EsriGrid, s *bufio.Scanner) {
	var err error
	var tmp float64

	//Getting right hand side of col - should validate but EsriGrids are standardised somewhat
	skipAndScan(s, 1)
	eg.ncols, err = strconv.Atoi(s.Text())
	Check(err)

	skipAndScan(s, 1)
	eg.nrows, err = strconv.Atoi(s.Text())
	Check(err)

	//TODO: There are xllcenter and yllcenter in some esri grids
	skipAndScan(s, 1)
	tmp, err = strconv.ParseFloat(s.Text(), 64)
	eg.xllcorner = int(tmp)
	Check(err)

	skipAndScan(s, 1)
	tmp, err = strconv.ParseFloat(s.Text(), 64)
	eg.yllcorner = int(tmp)
	Check(err)

	skipAndScan(s, 1)
	tmp, err = strconv.ParseFloat(s.Text(), 64)
	eg.cellsize = int(tmp)
	Check(err)

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
