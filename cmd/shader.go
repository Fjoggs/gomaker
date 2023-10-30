package gomaker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Shader struct {
	name     string
	lines    []string
	textures map[string]int
}

func extractTexturesFromUsedShaders(shadersFromMapFile map[string]int, shaderFolderPath string) (map[string]int, []string, []string) {
	shaderFiles := []string{}
	textures := map[string]int{}
	shaderNames := []string{}

	fsPath := addTrailingSlash(shaderFolderPath)
	directory, err := os.ReadDir(fsPath)

	for _, file := range directory {
		shaderFileName := file.Name()
		shaders := parseShaderFile(shadersFromMapFile, shaderFileName, shaderFolderPath)
		for _, shader := range shaders {
			delete(shadersFromMapFile, shader.name)
		}
		if len(shaders) > 0 {
			shaderFiles = append(shaderFiles, shaderFileName)
			textures, shaderNames = combineTexturesFromShaders(shaders, textures, shaderNames)
		}
	}

	if err != nil {
		log.Fatalf("Failed opening directory %v with path %s, error %s", directory, fsPath, err)
	}
	for key, value := range shadersFromMapFile {
		textures[key] = value
	}
	return textures, shaderNames, shaderFiles
}

func combineTexturesFromShaders(shaders []Shader, textures map[string]int, shaderNames []string) (map[string]int, []string) {
	for _, shader := range shaders {
		shaderNames = append(shaderNames, shader.name)
		for texture, count := range shader.textures {
			textures[texture] = count
		}
	}
	return textures, shaderNames
}

func parseShaderFile(shadersFromMapFile map[string]int, shaderFileName string, shaderFolderPath string) []Shader {
	fsPath := addTrailingSlash(shaderFolderPath) + shaderFileName
	file, err := os.Open(fsPath)

	if err != nil {
		log.Fatalf("Failed opening file %v with path %s, error %s", file, shaderFileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	shaders := []Shader{}
	shader := Shader{"", []string{}, map[string]int{}}
	parsingShader := false
	brackets := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "qer_editorimage") {
			fmt.Printf("Editor image %s, ignoring", line)
			continue
		}
		texture := getMaterial(line)
		if parsingShader {
			shader.lines = append(shader.lines, line)
		}
		if len(texture) > 0 {
			isShaderName := isShaderName(line)
			if isShaderName && shaderIsUsed(shadersFromMapFile, texture) {
				parsingShader = true
				shader.name = texture
				shader.lines = append(shader.lines, line)
			} else if parsingShader {
				shader.textures[texture] = shader.textures[texture] + 1
			}
		}
		if strings.Contains(line, "{") {
			brackets++
		} else if strings.Contains(line, "}") {
			brackets--
			if brackets == 0 && parsingShader {
				shaders = append(shaders, shader)
				shader = Shader{"", []string{}, map[string]int{}}
				parsingShader = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return shaders
}

func shaderIsUsed(shadersFromMapFile map[string]int, shaderName string) bool {
	_, ok := shadersFromMapFile[shaderName]
	return ok
}

func isShaderName(line string) bool {
	line = strings.Replace(line, "{", "", 1)
	line = strings.TrimSpace(line)
	hasWhitespaceRegexp := regexp.MustCompile(`\s`)
	hasWhitespace := hasWhitespaceRegexp.FindString(line)
	return len(hasWhitespace) == 0
}
