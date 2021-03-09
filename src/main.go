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

	allASCIIFiles := GetAllASCIIFiles(root)
	esriGrids := GenerateEsriGrids(allASCIIFiles)

	highY, lowY, highX, lowX, nRows, nCols, cellSize := GetEsriGridHighLow(esriGrids)

	// working but lts try reverse //this wont work if nrows, ncols, xllcorner, yllcorner arent integers
	//grid has to be square so see which one is bigger.
	stitchedGridLength := 0
	if (highY+nRows*cellSize)-lowY > (highX+nCols*cellSize)-lowX {
		fmt.Println((highY + nRows*cellSize) - lowY)
		stitchedGridLength = (highY + nRows*cellSize) - lowY
	} else {
		stitchedGridLength = (highX - lowX) + nCols*cellSize
	}
	stitchedGridLength = 17280
	fmt.Println("StitchedGridLength: ", stitchedGridLength)

	stitchedGrid := make([][]float32, stitchedGridLength)
	for i := range stitchedGrid {
		stitchedGrid[i] = make([]float32, stitchedGridLength)
	}

	highest := float32(math.MaxFloat32 * -1)
	lowest := float32(0) //float32(math.MaxFloat32)

	// maxX := highX + nRows*cellSize - lowX
	// maxY := highY + nRows*cellSize - lowY

	for _, eg := range esriGrids {
		//fmt.Println(eg.xllcorner, eg.yllcorner)
		fmt.Println(len(eg.grid))
		for y := 0; y < eg.nrows; y += cellSize {
			for x := 0; x < eg.ncols; x += cellSize {
				// fmt.Println(y+eg.yllcorner-lowY, " ", x+eg.xllcorner-lowX)
				val := eg.grid[nRows-1-y][x]
				if val == eg.noDataValue {
					val = -40 //sea level is 40m in C:S
				} else if val > highest {
					highest = val
				}

				// working but tiled weirdly stitchedGrid[y+eg.yllcorner-lowY][x+eg.xllcorner-lowX] = eg.grid[y][x]
				// working but inverted on x axis stitchedGrid[highY-eg.yllcorner+y][x+eg.xllcorner-lowX] = eg.grid[y][x]
				stitchedGrid[eg.yllcorner-lowY+y][x+eg.xllcorner-lowX] = val + 40
			}

		}

	}

	fmt.Println("highest: + 40 for map val", int(highest))
	fmt.Println("lowest", int(lowest))

	const mapResolution = 1081
	var mapScale = len(stitchedGrid) / mapResolution
	var colorScale = float32(64) //(float32(65535) / (highest - lowest))
	// fmt.Println(mapResolution, " ", len(stitchedGrid[0]))
	fmt.Println("Scale - map: ", mapScale, "colour: ", colorScale)
	fmt.Println(colorScale * (lowest + 40))

	//create image
	img := image.NewGray16(image.Rectangle{image.Point{0, 0}, image.Point{mapResolution, mapResolution}})
	// fmt.Println(stitchedGrid)
	for i := 0; i < mapResolution; i++ {
		for j := 0; j < mapResolution; j++ {
			img.Set(j, i, color.Gray16{uint16((stitchedGrid[j*mapScale][i*mapScale]) * colorScale)})

		}

	}

	f, _ := os.Create("image.png")
	png.Encode(f, img)
}

//scale it. Map is 17.28km^2

//  /1081 = ~16m - 15.9851988899
// its a series of intersections not the cell height the heightmap is though so maybe it wont be an issue
//https://community.simtropolis.com/forums/topic/72383-amis-berlin-assets-la-westside-map-released/?page=10

//rotate heightmap 180 deg since cities skylines sun rises from west to east
