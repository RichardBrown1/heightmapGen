package main

import (
	"fmt"
	"math"
	"os"
	"sort"
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

	fmt.Println(root[:1])

	if root[:1] == "/" {
		root = string(root[len(root)-1:])
	}

	fmt.Println("starting program at root path: " + root)

	//find ascii files
	allASCIIFiles := GetAllASCIIFiles(root)

	//fmt.Println(allASCIIFiles)
	esriGrids := GenerateEsriGrids(allASCIIFiles)
	//_ = esriGrids

	mapSizeX := (esriGrids[0].cellsize * float32(esriGrids[0].ncols))
	fmt.Println(mapSizeX)

	//todo: stitch heightmaps together
	//for now assume cellsize is the same and its all from 1 data source
	lowx := math.MaxFloat64
	lowy := math.MaxFloat64
	highx := math.MaxFloat64 * -1 //minFloat64
	highy := math.MaxFloat64 * -1
	cellSize := esriGrids[0].cellsize
	nRows := esriGrids[0].nrows
	nCols := esriGrids[0].ncols
	for _, eg := range esriGrids {
		//find lowest xll and yll corner
		if eg.xllcorner < lowx {
			lowx = eg.xllcorner
		} else {
			if eg.xllcorner > highx {
				highx = eg.xllcorner
			}
		}
		if eg.yllcorner < lowy {
			lowy = eg.yllcorner
		} else {
			if eg.yllcorner > highy {
				highy = eg.yllcorner
			}
		}
		if eg.cellsize != cellSize || eg.nrows != nRows || eg.ncols != nCols {
			fmt.Println("cellsize, row or col counts are not consistent ... will terminate this program isnt good enough to handle that ")
			panic(1)
		}
	}

	fmt.Println(
		" lowx: ", lowx,
		"\n lowy: ", lowy,
		"\n highx: ", highx,
		"\n highy: ", highy)

	sort.SliceStable(esriGrids, func(i, j int) bool {
		xDiff := esriGrids[i].xllcorner - esriGrids[j].xllcorner
		if xDiff != 0 {
			return esriGrids[i].xllcorner < esriGrids[j].xllcorner
		}
		return esriGrids[i].yllcorner < esriGrids[j].yllcorner
	})

	//allocate StitchedGrid
	//this wont work if nrows, ncols, xllcorner, yllcorner arent integers
	stitchedGrid := make([][]float32, int(highx+(float64(float32(nCols)*cellSize))-lowx))
	for i := range stitchedGrid {
		stitchedGrid[i] = make([]float32, int(highy+(float64(float32(nCols)*cellSize))-lowy))
	}

	//down to up; left to right
	i := 0
	for egx := lowx; egx <= highx && i < len(esriGrids); egx += float64(float32(nCols) * cellSize) {
		for egy := lowy; egy <= highy && i < len(esriGrids); egy += float64(float32(nRows) * cellSize) {
			fmt.Println("iterations", egx, egy)
			if esriGrids[i].xllcorner == egx && esriGrids[i].yllcorner == egy {
				for y := esriGrids.nrows; y > 0; y-- { //down to up
					for x := 0; x < esriGrids.ncols; x++ {
						stitchedGrid[egx-lowx][egy-lowy] = esriGrids[i][x][y]
					}
				}
				i++
			} else {
				fmt.Println("tile skipped... ")

			}
		}
	}
	//scale it. Map is 17.28km^2

	//  /1081 = ~16m - 15.9851988899
	// its a series of intersections not the cell height the heightmap is though so maybe it wont be an issue
	//https://community.simtropolis.com/forums/topic/72383-amis-berlin-assets-la-westside-map-released/?page=10

	//rotate heightmap 180 deg since cities skylines sun rises from west to east
}
