package entity

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gomaker/internal/material"
)

func IsEntity(line string) bool {
	match, err := regexp.MatchString("// entity", strings.ToLower(line))
	if err != nil {
		fmt.Println("Something went wrong", err)
	}

	return match
}

func ParseEntity(lines []string) map[string]int {
	textures := map[string]int{}
	modelPathLine := ""
	isModel := false
	for _, line := range lines {
		if !isModel {
			// Only overwrite variable if we haven't determined if it's a model yet
			isModel = strings.Contains(line, "misc_model")
		}
		if strings.Contains(line, "_remap") {
			texture := RemapTexture(line)
			textures[texture] = textures[texture] + 1
			return textures
		} else if strings.Contains(line, ".ase") {
			modelPathLine = line
		} else if strings.Contains(line, ".obj") {
			modelPathLine = strings.Replace(line, ".obj", ".mtl", 1)
		}
	}
	if isModel {
		modelPath := ModelPath(modelPathLine)
		textures = ParseModel(modelPath)
	}
	return textures
}

func ModelPath(line string) string {
	_, after, didCut := strings.Cut(line, "model")
	if didCut {
		after = strings.Replace(after, `"`, "", 3)
		return strings.TrimSpace(after)
	}
	return ""
}

func ParseModel(modelPath string) map[string]int {
	textures := map[string]int{}
	file, err := os.Open(modelPath)
	if err != nil {
		fmt.Printf("Failed opening file %v with path %s, error %s", file, modelPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	texture := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(modelPath, ".mtl") {
			texture = ObjTexture(line)
		} else {
			texture = AseTexture(line)
		}
		if len(texture) > 0 {
			textures[texture] = textures[texture] + 1
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return textures
}

func ObjTexture(line string) string {
	materialRegex := regexp.MustCompile("map_Kd")
	mat := materialRegex.FindString(line)
	if len(mat) > 0 {
		line = strings.TrimSpace(line)
		_, after, didCut := strings.Cut(line, `textures`)
		if didCut {
			texture := strings.ReplaceAll(after, "\\", "/")
			texture = strings.Replace(texture, `"`, "", 1)
			return material.GetMaterial(texture)
		}
	}
	return ""
}

func AseTexture(line string) string {
	materialRegex := regexp.MustCompile(`\*BITMAP[^_]`)
	mat := materialRegex.FindString(line)
	if len(mat) > 0 {
		line = strings.TrimSpace(line)
		_, after, didCut := strings.Cut(line, `textures`)
		if didCut {
			texture := strings.ReplaceAll(after, "\\", "/")
			texture = strings.Replace(texture, `"`, "", 1)
			return material.GetMaterial(texture)
		}
	}
	return ""
}

func RemapTexture(line string) string {
	return material.GetMaterial(line)
}
