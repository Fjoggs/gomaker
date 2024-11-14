package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gomaker/internal/brush"
	"gomaker/internal/entity"
	"gomaker/internal/material"
	"gomaker/internal/shader"
	"gomaker/internal/sound"
)

var (
	entityLines   []string
	parsingEntity bool
)

func init() {
	parsingEntity = false
}

func ReadMap(
	mapName string,
	baseFolderPath string,
) (map[string]int, map[string]int, []string, []string) {
	materials := map[string]int{}
	sounds := map[string]int{}
	file, err := os.Open(material.AddTrailingSlash(baseFolderPath) + "maps/" + mapName + ".map")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		AddMaterials(line, materials)
		sound.AddSounds(line, sounds)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	textures, shaderNames, shaderFiles := shader.ExtractTexturesFromUsedShaders(
		materials,
		fmt.Sprintf("%sscripts", material.AddTrailingSlash(baseFolderPath)),
	)

	textures = material.AddTexturePathWithExtension(
		textures,
		material.AddTrailingSlash(baseFolderPath),
	)

	return textures, sounds, shaderNames, shaderFiles
}

func AddMaterials(line string, materials map[string]int) {
	newMaterials := GetMaterials(line)
	if len(newMaterials) > 0 {
		MergeMaps(newMaterials, materials)
	}
}

func GetMaterials(line string) map[string]int {
	materials := map[string]int{}
	if entity.IsEntity(line) {
		parsingEntity = true
	} else if brush.IsBrush(line) {
		parsingEntity = false
	} else if IsClosingBracket(line) {
		if len(entityLines) > 0 {
			materials = HandleEntity(entityLines)
			entityLines = []string{}
		}
	} else {
		if parsingEntity {
			entityLines = append(entityLines, line)
		} else {
			texture := material.GetMaterial(line)
			if len(texture) > 0 {
				materials[texture] = materials[texture] + 1
			}
		}
	}
	return materials
}

func HandleBrush(line string) string {
	parsingEntity = false
	return material.GetMaterial(line)
}

func HandleEntity(lines []string) map[string]int {
	parsingEntity = false
	return entity.ParseEntity(lines)
}

func IsClosingBracket(line string) bool {
	return strings.Contains(line, "}")
}

func MergeMaps(source map[string]int, destination map[string]int) {
	for key, count := range source {
		destination[key] = count
	}
}
