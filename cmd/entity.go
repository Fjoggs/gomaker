package gomaker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func isEntity(line string) bool {
	match, err := regexp.MatchString("// entity", strings.ToLower(line))

	if err != nil {
		log.Fatal("Something went wrong", err)
	}

	return match
}

func parseEntity(lines []string) []string {
	textures := []string{}
	modelPathLine := ""
	hasRemap := false
	for _, line := range lines {
		if strings.Contains(line, "_remap") {
			hasRemap = true
			textures = append(textures, remapTexture(line))
			return textures
		} else if strings.Contains(line, ".ase") {
			modelPathLine = line
			fmt.Println("line", line)
		} else if strings.Contains(line, ".obj") {
			fmt.Println("line", line)
			modelPathLine = strings.Replace(line, ".obj", ".mtl", 1)
		}
	}
	if !hasRemap {
		modelPath := modelPath(modelPathLine)
		textures = parseModel(modelPath)
	}
	return textures
}

func modelPath(line string) string {
	_, after, didCut := strings.Cut(line, "model")
	if didCut {
		return strings.TrimSpace(after)
	}
	return ""
}

func parseModel(modelPath string) []string {
	textures := []string{}
	file, err := os.Open(modelPath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	texture := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(modelPath, ".mtl") {
			texture = objTexture(line)
		} else {
			texture = aseTexture(line)
		}
		if len(texture) > 0 {
			textures = append(textures, texture)
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return textures
}

func objTexture(line string) string {
	materialRegex := regexp.MustCompile("map_Kd")
	material := materialRegex.FindString(line)
	if len(material) > 0 {
		line = strings.TrimSpace(line)
		_, after, didCut := strings.Cut(line, `textures`)
		if didCut {
			texture := strings.ReplaceAll(after, "\\", "/")
			texture = strings.Replace(texture, `"`, "", 1)
			return isTexture(texture)
		}
	}
	return ""
}

func aseTexture(line string) string {
	materialRegex := regexp.MustCompile(`\*BITMAP[^_]`)
	material := materialRegex.FindString(line)
	if len(material) > 0 {
		line = strings.TrimSpace(line)
		_, after, didCut := strings.Cut(line, `textures`)
		if didCut {
			texture := strings.ReplaceAll(after, "\\", "/")
			texture = strings.Replace(texture, `"`, "", 1)
			return isTexture(texture)
		}
	}
	return ""
}

func remapTexture(line string) string {
	return isTexture(line)
}
