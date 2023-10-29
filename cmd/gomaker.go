package gomaker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	filePath := "resources/test.map"
	textures := readMap(filePath, "textures/")
	fmt.Println(textures)
}

func readMap(path string, textureFolderPath string) Materials {
	textures := []string{}
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newMaterials := getMaterials(scanner.Text())
		textures = combineSlices([][]string{textures, newMaterials})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	materials := sortMaterials(textures, textureFolderPath)

	return materials
}

func getMaterials(line string) []string {
	materials := []string{}
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
				materials = append(materials, texture)
			}
		}
	}
	return materials
}

func handleBrush(line string) string {
	parsingEntity = false
	return getMaterial(line)
}

func handleEntity(lines []string) []string {
	parsingEntity = false
	return parseEntity(lines)
}

func isClosingBracket(line string) bool {
	return strings.Contains(line, "}")
}

func combineSlices(slices [][]string) []string {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]string, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}

func sortMaterials(materials []string, basePath string) Materials {
	sorted := Materials{make(map[string]int), make(map[string]int)}
	for _, material := range materials {
		isT, filePath := isTexture(material, basePath)
		if isT {
			sorted.textures[filePath] = sorted.textures[filePath] + 1
			// It can also be a shader
			sorted.shaders[material] = sorted.shaders[material] + 1
		} else {
			sorted.shaders[material] = sorted.shaders[material] + 1
		}
	}
	return sorted
}
