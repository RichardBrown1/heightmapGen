package main

import (
	"fmt"
	"math"
	"os"
)

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

	mapSizeX := (esriGrids[0].cellsize * esriGrids[0].ncols)
	fmt.Println(mapSizeX)

	//todo: stitch heightmaps together
	//for now assume cellsize is the same and its all from 1 data source
	// lowX := math.MaxFloat64
	// lowY := math.MaxFloat64
	// highX := math.MaxFloat64 * -1 //minFloat64
	// highY := math.MaxFloat64 * -1
	lowX := math.MaxInt32
	lowY := math.MaxInt32
	highX := math.MaxInt32 * -1 //minFloat64
	highY := math.MaxInt32 * -1

	cellSize := esriGrids[0].cellsize
	nRows := esriGrids[0].nrows
	nCols := esriGrids[0].ncols
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
		"\n highY: ", highY)

	// sort.SliceStable(esriGrids, func(i, j int) bool {
	// 	if esriGrids[i].xllcorner == esriGrids[j].xllcorner {
	// 		return esriGrids[i].yllcorner < esriGrids[j].yllcorner
	// 	}
	// 	return esriGrids[i].xllcorner < esriGrids[j].xllcorner
	// })

	//this wont work if nrows, ncols, xllcorner, yllcorner arent integers
	stitchedGrid := make([][]float32, int((highX+nRows*cellSize)-lowX))
	for i := range stitchedGrid {
		stitchedGrid[i] = make([]float32, int((highY+nCols*cellSize)-lowY))
	}

	for y := highY; y <= lowY; y += cellSize {
		for x := highX; x <= lowX; x += cellSize {
			for _, eg := range esriGrids {
				if eg.xllcorner == x && eg.yllcorner == y {
					for egY := 0; egY < eg.nrows; egY++ {
						for egX := 0; egX < eg.nrows; egX++ {
							stitchedGrid[highX-x][highY-y] = eg.grid[egY][egX]
						}
					}

				}
			}
		}
	}

	//allocate finalgrid.
	const mapResolution = 1081
	finalGrid := make([][]float32, mapResolution)
	for i := range finalGrid {
		finalGrid[i] = make([]float32, mapResolution)
	}

	//create image
	// img := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{mapResolution, mapResolution}})

	// for y, eg := range finalGrid {
	// 	for x, eg := range finalGrid[i] {
	// 		// img.Set(x,y, color.Gray16{} )
	// 	}
	// }
}

//this wont work if nrows, ncols, xllcorner, yllcorner arent integers
// stitchedGrid := make([][]float32, int(highX+(float64(float32(nCols)*cellSize))-lowX))
// for i := range stitchedGrid {
// 	stitchedGrid[i] = make([]float32, int(highY+(float64(float32(nCols)*cellSize))-lowY))
// }

//todo remove floats and only allow integer cellsizes and w/e

//down to up; left to right
// i := 0
// for egx := lowX; egx <= highX && i < len(esriGrids); egx += float64(float32(nCols) * cellSize) {
// 	for egy := lowY; egy <= highY && i < len(esriGrids); egy += float64(float32(nRows) * cellSize) {
// 		fmt.Println("iterations", egx, egy)
// 		if esriGrids[i].xllcorner == egx && esriGrids[i].yllcorner == egy {
// 			for y := esriGrids.nrows; y > 0; y-- { //down to up
// 				for x := 0; x < esriGrids.ncols; x++ {
// 					stitchedGrid[egx-lowX][egy-lowY] = esriGrids[i][x][y]
// 				}
// 			}
// 			i++
// 		} else {
// 			fmt.Println("tile skipped... ")

// 		}
// 	}
// }
//scale it. Map is 17.28km^2

//  /1081 = ~16m - 15.9851988899
// its a series of intersections not the cell height the heightmap is though so maybe it wont be an issue
//https://community.simtropolis.com/forums/topic/72383-amis-berlin-assets-la-westside-map-released/?page=10

//rotate heightmap 180 deg since cities skylines sun rises from west to east
