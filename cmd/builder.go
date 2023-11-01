package gomaker

import (
	"fmt"
	"log"
	"os"
)

func createDirectory(mapName string) {
	pk3Folder, pk3Err := os.MkdirTemp("", "gomaker-pk3-"+mapName)
	if pk3Err != nil {
		log.Fatalf("It blew up mate %s", pk3Err)
	}

	textureFolder, _ := os.MkdirTemp(pk3Folder, "textures")
	fmt.Println(pk3Folder)
	fmt.Println(textureFolder)
	os.RemoveAll(pk3Folder)
	_, err := os.Stat(pk3Folder)
	if err != nil {
		fmt.Printf("Stat failed with error %s", err)
	}
}
