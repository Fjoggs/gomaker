package gomaker

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func createPk3(baseq3Folder string, resources []string, mapName string, overwrite bool) {
	if overwrite {
		deleteFolderAndSubFolders(
			fmt.Sprintf("%s%s", addTrailingSlash(baseq3Folder), addTrailingSlash(mapName)),
		)
	}

	createDirectory("output", "")
	for _, resource := range resources {
		fmt.Printf("Checking if resource %s exists\n", resource)
		addResourceIfExists(baseq3Folder, resource, "output")
	}

	err := zipOutputFolderAsPk3("output", mapName)
	if err != nil {
		fmt.Printf("Eyo? %s", err)
	}
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
		fmt.Printf("It blew up mate %s\n", err)
		return false
	}

	return true
}

func zipOutputFolderAsPk3(outputFolder string, mapName string) error {
	targetPath := addTrailingSlash(outputFolder) + mapName + ".pk3"
	file, err := os.Create(targetPath)
	if err != nil {
		fmt.Printf("Error occured while creating zip: %s", err)
		return err
	}

	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	sourcePath := addTrailingSlash(outputFolder)
	return filepath.WalkDir(sourcePath, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fileInfo, err := dir.Info()
		if err != nil {
			return err
		}

		if dir.Name() == mapName {
			fmt.Println("Omitting map name folder")
			return nil
		}

		if dir.Name() == mapName+".pk3" {
			fmt.Println("Omitting itself (wut)")
			return nil
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		name := strings.Replace(path, sourcePath, "", 1)
		header.Method = zip.Deflate
		header.Name = name
		if dir.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(headerWriter, file)
		return err
	})
}

func addResourceIfExists(baseq3Folder string, resourcePath string, outputFolder string) string {
	path := fmt.Sprintf("%s%s", addTrailingSlash(baseq3Folder), resourcePath)
	fmt.Printf("path %s\n", path)

	_, exists := os.Stat(path)
	if exists != nil {
		fmt.Printf("Resource does not exist: %s\n", path)
		return ""
	}

	sourceFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Something went wrong opening source file: %s\n", err)
		return ""
	}

	destPath := addTrailingSlash(outputFolder) + resourcePath
	destFolder := extractFolderPaths(destPath)
	_, exists = os.Stat(destFolder)
	if exists != nil {
		fmt.Printf("Destination path does not exist: %s\n", destPath)
		err = os.MkdirAll(destFolder, 0777)
		if err != nil {
			fmt.Printf("MkdirAll returned error: %s", err)
			return ""
		}
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Something went wrong creating target file: %s\n", err)
		return ""
	}

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		fmt.Printf("Something went wrong while copying: %s", err)
		return ""
	}

	fmt.Printf("Added resource %s\n", destPath)
	return destPath
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

func getCfgFile(baseq3Folder string, mapName string) string {
	mapFilePath := fmt.Sprintf("%scfg-maps/%s.cfg", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(mapFilePath)
	if err != nil {
		fmt.Printf("No .cfg file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found .cfg file %s.cfg\n", mapName)
	return fmt.Sprintf("cfg-maps/%s.cfg", mapName)
}

func getReadme(baseq3Folder string, mapName string) string {
	mapFilePath := fmt.Sprintf("%s%s.txt", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(mapFilePath)
	if err != nil {
		fmt.Printf("No .txt file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found .txt file %s.txt\n", mapName)
	return fmt.Sprintf("%s.txt", mapName)
}

func getBspFile(baseq3Folder string, mapName string) string {
	mapFilePath := fmt.Sprintf("%smaps/%s.bsp", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(mapFilePath)
	if err != nil {
		fmt.Printf("No .bsp file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found .bsp file %s.bsp\n", mapName)
	return fmt.Sprintf("maps/%s.bsp", mapName)
}

func getMapFile(baseq3Folder string, mapName string) string {
	mapFilePath := fmt.Sprintf("%smaps/%s.map", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(mapFilePath)
	if err != nil {
		fmt.Printf("No .map file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found .map file %s.map\n", mapName)
	return fmt.Sprintf("maps/%s.map", mapName)
}

func getExternalLightmaps(baseq3Folder string, mapName string) []string {
	mapFilePath := fmt.Sprintf("%smaps/%s/", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(mapFilePath)
	lightmaps := []string{}
	if err != nil {
		fmt.Printf("No external lightmap folder found: %s\n", err)
		return lightmaps
	}

	err = filepath.WalkDir(mapFilePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if mapName == d.Name() {
			fmt.Printf("Ignoring root folder %s", d.Name())
			return err
		}

		lightmaps = append(lightmaps, fmt.Sprintf("maps/%s/%s", mapName, d.Name()))
		return err
	})
	if err != nil {
		fmt.Printf("an error occured while walking %s\n", err)
		return lightmaps
	}

	defer file.Close()
	fmt.Printf("Found lightmaps %s\n", lightmaps)
	return lightmaps
}

func getArenaFile(baseq3Folder string, mapName string) string {
	arenaFilePath := fmt.Sprintf("%sscripts/%s.arena", addTrailingSlash(baseq3Folder), mapName)
	file, err := os.Open(arenaFilePath)
	if err != nil {
		fmt.Printf("No arena file found: %s", err)
		return ""
	}

	defer file.Close()
	fmt.Printf("Found arena file %s.arena\n", mapName)
	return fmt.Sprintf("scripts/%s.arena", mapName)
}

func getLevelshot(baseq3Folder string, mapName string) string {
	levelshotsPath := fmt.Sprintf("%slevelshots/%s", addTrailingSlash(baseq3Folder), mapName)
	jpg := fmt.Sprintf("%s.jpg", levelshotsPath)
	tga := fmt.Sprintf("%s.tga", levelshotsPath)
	jpgFile, jpgErr := os.Open(jpg)

	if jpgErr == nil {
		return fmt.Sprintf("levelshots/%s.jpg", mapName)
	}
	defer jpgFile.Close()

	tgaFile, tgaErr := os.Open(tga)

	if tgaErr == nil {
		return fmt.Sprintf("levelshots/%s.tga", mapName)
	}

	defer tgaFile.Close()

	return ""
}

func extractFolderPaths(fullPath string) string {
	if strings.Contains(fullPath, ".") {
		split := strings.Split(fullPath, "/")
		if len(split) > 0 {
			split = split[:len(split)-1]
			fullPath = strings.Join(split, "/")
		}
		return fullPath
	} else {
		return fullPath
	}
}
