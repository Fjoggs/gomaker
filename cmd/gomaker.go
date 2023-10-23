package gomaker

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	filePath := "resources/test.map"
	readMap(filePath)
}

func readMap(path string) []string {
	textures := []string{}
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	parsingEntity := false
	entityLines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if parsingEntity {
			if isEntity(line) {
				// New entity
				parsingEntity = true
				entityLines = []string{}
			} else if isBrush(line) {
				// Something
				parsingEntity = false
			} else {
				entityLines = append(entityLines, line)
			}
		}
		if isEntity(line) {
			parsingEntity = true
		}
		texture := readLine(line)
		if len(texture) > 0 {
			fmt.Println("texture", texture)
			textures = append(textures, texture)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return textures
}

func readLine(line string) string {
	return isTexture(line)
}
