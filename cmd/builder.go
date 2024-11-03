package gomaker

import (
	"fmt"
	"io"
	"log"
	"os"
)

func createPk3(outputFolder string, mapName string, overwrite bool) bool {
	if overwrite {
		deleteFolderAndSubFolders(fmt.Sprintf("%s%s", addTrailingSlash(outputFolder), addTrailingSlash(mapName)))
	}

	createDirectory("output", "")
	pk3Dir := createDirectory(mapName, "output")
	if pk3Dir {
		pk3DirPath := fmt.Sprintf("%s%s", addTrailingSlash(outputFolder), addTrailingSlash(mapName))
		envDir := createDirectory("env", pk3DirPath)
		mapsDir := createDirectory("maps", pk3DirPath)
		textureDir := createDirectory("textures", pk3DirPath)
		scriptsDir := createDirectory("scripts", pk3DirPath)
		soundsDir := createDirectory("sounds", pk3DirPath)
		levelshotsDir := createDirectory("levelshots", pk3DirPath)

		fmt.Println("Getting arenafile")
		arenaFile := getArenaFile(mapName)
		fmt.Println("Getting levelshot")
		levelshot := getLevelshot(mapName)

		fmt.Println("Adding arenafile")
		addResourceIfExists(arenaFile, mapName, outputFolder)
		fmt.Println("Adding levelshot")
		addResourceIfExists(levelshot, mapName, outputFolder)

		fmt.Println(pk3Dir)
		fmt.Println(envDir)
		fmt.Println(mapsDir)
		fmt.Println(textureDir)
		fmt.Println(scriptsDir)
		fmt.Println(soundsDir)
		fmt.Println(levelshotsDir)
		return true
	}
	return false
}

func createDirectory(folderName string, mapName string) bool {
	path := addTrailingSlash(mapName) + folderName
	_, existErr := os.Stat(path)

	if existErr == nil {
		fmt.Printf("Removing existing folder %s and contents\n", path)
		deleteFolderAndSubFolders(path)
	}

	err := os.Mkdir(path, 0777)
	if err != nil {
		log.Fatalf("It blew up mate %s\n", err)
		return false
	}

	return true
}

func addResourceIfExists(resourceName string, mapName string, targetFolder string) bool {
	path := addTrailingSlash(mapName) + resourceName

	_, exists := os.Stat(path)
	if exists != nil {
		fmt.Printf("Resource does not exist: %s\n", path)
		return false
	}

	sourceFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Something went wrong opening source file: %s\n", err)
		return false
	}

	destPath := addTrailingSlash(targetFolder) + resourceName
	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Something went wrong creating target file: %s\n", err)
		return false
	}

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Printf("Something went wrong while copying: %s", err)
		return false
	}

	fmt.Printf("Added resource %s\n", resourceName)
	return true
}

func deleteFolderAndSubFolders(folder string) {
	path := addTrailingSlash(folder)

	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Stat failed with error %s\n", err)
	}

	fmt.Printf("Removing %s\n", path)
	os.RemoveAll(path)
}

func getArenaFile(mapName string) string {
	arenaFilePath := fmt.Sprintf("%sscripts/%s.arena", addTrailingSlash(mapName), mapName)
	file, err := os.Open(arenaFilePath)

	if err != nil {
		fmt.Printf("No arena file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found arena file %s.arena\n", mapName)
	return fmt.Sprintf("scripts/%s.arena", mapName)
}

func getLevelshot(mapName string) string {
	screenshotPath := fmt.Sprintf("%slevelshots/%s", addTrailingSlash(mapName), mapName)
	jpg := fmt.Sprintf("%s.jpg", screenshotPath)
	tga := fmt.Sprintf("%s.tga", screenshotPath)
	jpgFile, jpgErr := os.Open(jpg)

	if jpgErr == nil {
		return fmt.Sprintf("levelshots/%s.jpg", mapName)
	} else {
		log.Printf("Failed opening jpg file with path %s, error %s\n", screenshotPath, jpgErr)
	}
	defer jpgFile.Close()

	tgaFile, tgaErr := os.Open(tga)

	if tgaErr == nil {
		return fmt.Sprintf("levelshots/%s.tga", mapName)
	} else {
		log.Printf("Failed opening tga file with path %s, error %s\n", screenshotPath, tgaErr)
	}

	defer tgaFile.Close()

	return ""
}
