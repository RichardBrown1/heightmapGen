package main

import (
	"fmt"
	"os"
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
	_ = esriGrids

	mapSizeX := (esriGrids[0].cellsize * float32(esriGrids[0].ncols))
	fmt.Println(mapSizeX)

	//todo: stitch heightmaps together

	//scale it. Map is 17.28km^2

	//  /1081 = ~16m - 15.9851988899
	// its a series of intersections not the cell height the heightmap is though so maybe it wont be an issue
	//https://community.simtropolis.com/forums/topic/72383-amis-berlin-assets-la-westside-map-released/?page=10

	//rotate heightmap 180 deg since cities skylines sun rises from west to east
}
