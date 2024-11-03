package gomaker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Materials struct {
	textures map[string]int
	shaders  map[string]int
}

var entityLines []string
var parsingEntity bool

func init() {
	parsingEntity = false
}

func main() {
	start := time.Now()
	mapName := "test.map"
	textures, sounds, shaderNames := readMap(mapName, "resources")
	arenaFile := getArenaFile("testmap")
	fmt.Println(textures)
	fmt.Println(sounds)
	fmt.Println(shaderNames)
	fmt.Println(arenaFile)
	elapsed := time.Since(start)
	fmt.Println("Elapsed time", elapsed)
}

func readMap(mapName string, baseFolderPath string) (map[string]int, map[string]int, []string) {
	materials := map[string]int{}
	sounds := map[string]int{}
	file, err := os.Open(addTrailingSlash(baseFolderPath) + "maps/" + mapName)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		addMaterials(line, materials)
		addSounds(line, sounds)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	textures, shaderNames, _ := extractTexturesFromUsedShaders(materials, "resources/scripts")

	textures = addTextureFileExtension(textures, addTrailingSlash(baseFolderPath))

	return textures, sounds, shaderNames
}

func addMaterials(line string, materials map[string]int) {
	newMaterials := getMaterials(line)
	if len(newMaterials) > 0 {
		mergeMaps(newMaterials, materials)
	}
}

func addSounds(line string, sounds map[string]int) {
	sound := getSound(line)
	if len(sound) > 0 {
		sounds[sound] = sounds[sound] + 1
	}
}

func getMaterials(line string) map[string]int {
	materials := map[string]int{}
	if isEntity(line) {
		parsingEntity = true
	} else if isBrush(line) {
		parsingEntity = false
	} else if isClosingBracket(line) {
		if len(entityLines) > 0 {
			materials = handleEntity(entityLines)
			entityLines = []string{}
		}
	} else {
		if parsingEntity {
			entityLines = append(entityLines, line)
		} else {
			texture := getMaterial(line)
			if len(texture) > 0 {
				materials[texture] = materials[texture] + 1
			}
		}
	}
	return materials
}

func handleBrush(line string) string {
	parsingEntity = false
	return getMaterial(line)
}

func handleEntity(lines []string) map[string]int {
	parsingEntity = false
	return parseEntity(lines)
}

func isClosingBracket(line string) bool {
	return strings.Contains(line, "}")
}

func mergeMaps(source map[string]int, destination map[string]int) {
	for key, count := range source {
		destination[key] = count
	}
}
