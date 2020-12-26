package heightmapgen

import (
	"fmt"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	root := os.Args[1]

	fmt.Println("starting program at root path: " + root)

	rootDir, err := listSubDirs(root)
	check(err)
	fmt.Println(rootDir)
}
