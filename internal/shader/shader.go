package shader

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gomaker/internal/material"
)

type Shader struct {
	Name     string
	Lines    []string
	Textures map[string]int
}

func ExtractTexturesFromUsedShaders(
	shadersFromMapFile map[string]int,
	shaderFolderPath string,
) (map[string]int, []string, []string) {
	shaderFiles := []string{}
	textures := map[string]int{}
	shaderNames := []string{}

	fsPath := material.AddTrailingSlash(shaderFolderPath)
	directory, err := os.ReadDir(fsPath)

	for _, file := range directory {
		shaderFileName := file.Name()
		shaders := ParseShaderFile(shadersFromMapFile, shaderFileName, shaderFolderPath)
		for _, shader := range shaders {
			delete(shadersFromMapFile, shader.Name)
		}
		if len(shaders) > 0 {
			shaderFiles = append(shaderFiles, shaderFileName)
			textures, shaderNames = CombineTexturesFromShaders(shaders, textures, shaderNames)
		}
	}

	if err != nil {
		fmt.Printf("Failed opening directory %v with path %s, error %s", directory, fsPath, err)
	}
	for key, value := range shadersFromMapFile {
		textures[key] = value
	}
	return textures, shaderNames, shaderFiles
}

func CombineTexturesFromShaders(
	shaders []Shader,
	textures map[string]int,
	shaderNames []string,
) (map[string]int, []string) {
	for _, shader := range shaders {
		shaderNames = append(shaderNames, shader.Name)
		for texture, count := range shader.Textures {
			textures[texture] = count
		}
	}
	return textures, shaderNames
}

func ParseShaderFile(
	shadersFromMapFile map[string]int,
	shaderFileName string,
	shaderFolderPath string,
) []Shader {
	fsPath := material.AddTrailingSlash(shaderFolderPath) + shaderFileName
	file, err := os.Open(fsPath)
	if err != nil {
		fmt.Printf("Failed opening file %v with path %s, error %s", file, shaderFileName, err)
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
			continue
		}
		texture := material.GetMaterial(line)
		if parsingShader {
			shader.Lines = append(shader.Lines, line)
		}
		if len(texture) > 0 {
			isShaderName := IsShaderName(line)
			if isShaderName && ShaderIsUsed(shadersFromMapFile, texture) {
				parsingShader = true
				shader.Name = texture
				shader.Lines = append(shader.Lines, line)
			} else if parsingShader {
				shader.Textures[texture] = shader.Textures[texture] + 1
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
		fmt.Println(err)
	}
	return shaders
}

func ShaderIsUsed(shadersFromMapFile map[string]int, shaderName string) bool {
	_, ok := shadersFromMapFile[shaderName]
	return ok
}

func IsShaderName(line string) bool {
	line = strings.Replace(line, "{", "", 1)
	line = strings.TrimSpace(line)
	hasWhitespaceRegexp := regexp.MustCompile(`\s`)
	hasWhitespace := hasWhitespaceRegexp.FindString(line)
	return len(hasWhitespace) == 0
}
