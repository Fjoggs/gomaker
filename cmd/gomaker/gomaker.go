package main

import (
	"fmt"
	"os"
	"time"

	"gomaker/internal/builder"
)

func main() {
	start := time.Now()

	if len(os.Args[1:]) > 1 {
		mapName := os.Args[1]
		basePath := os.Args[2]
		pk3Path := builder.BuildPk3(mapName, basePath)
		fmt.Printf("Pk3 built and copied to %s\n", pk3Path)
	} else {
		mapName := os.Getenv("MAPNAME")
		basePath := os.Getenv("Q3_BASEPATH")
		if len(mapName) == 0 || len(basePath) == 0 {
			fmt.Println("Either pass map name and base path as arguments, or export env variables MAPNAME and Q3_BASEPATH")
		} else {
			pk3Path := builder.BuildPk3(mapName, basePath)
			fmt.Printf("Pk3 built and copied to %s\n", pk3Path)
		}
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed time", elapsed)
}
