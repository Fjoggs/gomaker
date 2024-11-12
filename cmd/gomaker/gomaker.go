package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"strings"
	"time"

	gomaker "gomaker/internal"
)

type Materials struct {
	textures map[string]int
	shaders  map[string]int
}

var (
	entityLines   []string
	parsingEntity bool
)

func init() {
	parsingEntity = false
}

func main() {
	start := time.Now()

	if len(os.Args[1:]) < 2 {
		fmt.Println("Eyo?")
		elapsed := time.Since(start)
		fmt.Println("Elapsed time", elapsed)
		return
	}

	mapName := os.Args[1]
	basePath := os.Args[2]

	resources := []string{}

	readme := getReadme(basePath, mapName)
	resources = append(resources, readme)

	cfg := getCfgFile(basePath, mapName)
	resources = append(resources, cfg)

	mapFile := getMapFile(basePath, mapName)
	resources = append(resources, mapFile)

	bspFile := getBspFile(basePath, mapName)
	resources = append(resources, bspFile)

	arenaFile := getArenaFile(basePath, mapName)
	resources = append(resources, arenaFile)

	levelshot := getLevelshot(basePath, mapName)
	resources = append(resources, levelshot)

	lightmaps := getExternalLightmaps(basePath, mapName)
	fmt.Printf("Lightmaps %v", lightmaps)

	textures, sounds, shaderNames, shaderFiles := readMap(mapName, basePath)

	for texture := range maps.Keys(textures) {
		resources = append(resources, texture)
	}

	for sound := range maps.Keys(sounds) {
		resources = append(resources, sound)
	}

	for _, shaderFile := range shaderFiles {
		resources = append(resources, "scripts/"+shaderFile)
	}

	resources = append(resources, lightmaps...)
	resources = append(resources, shaderNames...)

	gomaker.CreatePk3(basePath, resources, mapName, false)

	elapsed := time.Since(start)
	fmt.Println("Elapsed time", elapsed)
}

func readMap(
	mapName string,
	baseFolderPath string,
) (map[string]int, map[string]int, []string, []string) {
	materials := map[string]int{}
	sounds := map[string]int{}
	file, err := os.Open(addTrailingSlash(baseFolderPath) + "maps/" + mapName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		addMaterials(line, materials)
		addSounds(line, sounds)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	textures, shaderNames, shaderFiles := extractTexturesFromUsedShaders(
		materials,
		"resources/scripts",
	)

	textures = addTexturePathWithExtension(textures, addTrailingSlash(baseFolderPath))

	return textures, sounds, shaderNames, shaderFiles
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
