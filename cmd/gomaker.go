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
	textures, shaderNames := readMap(filePath, "textures/")
	fmt.Println(textures)
	fmt.Println(shaderNames)
}

func readMap(path string, textureFolderPath string) (map[string]int, []string) {
	unsortedMaterials := []string{}
	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		newMaterials := getMaterials(scanner.Text())
		unsortedMaterials = combineSlices([][]string{unsortedMaterials, newMaterials})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	materials := sortMaterials(unsortedMaterials, textureFolderPath)
	textures, shaderNames, _ := extractTexturesFromUsedShaders(materials.shaders, "resources/scripts")

	textures = addExtension(textures, "resources/textures")
	return textures, shaderNames
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

func addExtension(textures map[string]int, basePath string) map[string]int {
	returnValue := map[string]int{}
	for material := range textures {
		isT, filePath := isTexture(material, basePath)
		if isT {
			returnValue[filePath] = returnValue[filePath] + 1
		}
	}
	return returnValue
}
