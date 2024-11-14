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
		builder.MakePk3(mapName, basePath)
	}
	elapsed := time.Since(start)
	fmt.Println("Elapsed time", elapsed)
}
