package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
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
	// stitchedGrid := make([][]float32, int((highX+nRows*cellSize)-lowX))
	// for i := range stitchedGrid {
	// 	stitchedGrid[i] = make([]float32, int((highY+nCols*cellSize)-lowY))
	// }

	// working but lts try reverse //this wont work if nrows, ncols, xllcorner, yllcorner arent integers
	//grid has to be square so see which one is bigger.
	stitchedGridLength := 0
	if (highY+nRows*cellSize)-lowY > (highX+nCols*cellSize)-lowX {
		stitchedGridLength = (highY + nRows*cellSize) - lowY
	} else {
		stitchedGridLength = (highX + nCols*cellSize) - lowX
	}

	stitchedGrid := make([][]float32, stitchedGridLength)
	for i := range stitchedGrid {
		stitchedGrid[i] = make([]float32, stitchedGridLength)
	}

	// stitchedGrid := make([][]float32, int((highY+nRows*cellSize)-lowY))
	// for i := range stitchedGrid {
	// 	stitchedGrid[i] = make([]float32, int((highX+nCols*cellSize)-lowX))
	// }

	highest := float32(math.MaxFloat32 * -1)
	lowest := float32(-40) //float32(math.MaxFloat32)

	// maxX := highX + nRows*cellSize - lowX
	// maxY := highY + nRows*cellSize - lowY

	for _, eg := range esriGrids {
		fmt.Println(eg.xllcorner, eg.yllcorner)
		for y := 0; y < eg.nrows; y++ {
			for x := 0; x < eg.ncols; x++ {
				// fmt.Println(y+eg.yllcorner-lowY, " ", x+eg.xllcorner-lowX)
				if eg.grid[y][x] == eg.noDataValue {
					eg.grid[y][x] = -40
				} else if eg.grid[y][x] > highest {
					highest = eg.grid[y][x]
				}

				// working but tiled weirdly stitchedGrid[y+eg.yllcorner-lowY][x+eg.xllcorner-lowX] = eg.grid[y][x]
				// working but inverted on x axis stitchedGrid[highY-eg.yllcorner+y][x+eg.xllcorner-lowX] = eg.grid[y][x]
				stitchedGrid[eg.yllcorner-lowY+y][x+eg.xllcorner-lowX] = eg.grid[nRows-1-y][x]
			}

		}

	}

	fmt.Println("highest:", int(highest))
	fmt.Println("lowest", int(lowest))

	const mapResolution = 1081
	var mapScale = len(stitchedGrid) / mapResolution
	var colorScale = float32(64) //(float32(65535) / (highest - lowest))
	// fmt.Println(mapResolution, " ", len(stitchedGrid[0]))
	fmt.Println(mapScale, " ", colorScale)
	fmt.Println(colorScale * (lowest + 40))

	//create image
	img := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{mapResolution, mapResolution}})
	// fmt.Println(stitchedGrid)
	for i := 0; i < mapResolution; i++ {
		for j := 0; j < mapResolution; j++ {

			//fmt.Println(stitchedGrid[i][j] + 40)
			// img.Set(j, i, color.Gray{uint8(x + y)})
			img.Set(j, i, color.Gray16{uint16((stitchedGrid[j*mapScale][i*mapScale] + 40) * colorScale)})
			// img.Set(j, i, color.Gray{uint8(stitchedGrid[x][y] + 40)})

		}

	}

	// i := 0
	// for y := 0; y < mapResolution; y += mapScale {
	// 	j := 0
	// 	for x := 0; x < mapResolution; x += mapScale {
	// 		fmt.Println(mapResolution, " ", mapScale)
	// 		fmt.Println(stitchedGrid[x][y] + 40)
	// 		img.Set(j, i, color.Gray{uint8(x + y)})
	// 		// img.Set(j, i, color.Gray{uint8(stitchedGrid[x][y] + 40)})
	// 		j++
	// 	}
	// 	i++
	// }

	f, _ := os.Create("image.png")
	png.Encode(f, img)

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
