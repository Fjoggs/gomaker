package gomaker

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func getMaterial(line string) string {
	textureRegex := regexp.MustCompile(`((\w+[\/_-]*)+\/((\w)+[\/_-]*)*)+`)
	texture := textureRegex.FindString(line)
	if len(texture) > 0 {
		if isCustomMaterial(texture) {
			return formatPath(texture)
		}
	}
	return ""
}

func isCustomMaterial(texture string) bool {
	common := [5]string{"common/", "common_alphascale/", "sfx/", "liquids/", "effects/"}
	for _, value := range common {
		if strings.Contains(texture, value) {
			return false
		}
	}
	return true
}

func formatPath(texture string) string {
	return strings.Replace(texture, "textures/", "", 1)
}

func isTexture(material string, baseFolderPath string) (bool, string) {
	fsPath := addTrailingSlash(baseFolderPath) + "textures/" + material
	jpg := fmt.Sprintf("%s.jpg", fsPath)
	tga := fmt.Sprintf("%s.tga", fsPath)
	jpgFile, jpgErr := os.Open(jpg)

	if jpgErr == nil {
		return true, fmt.Sprintf("%s.jpg", material)
	} else {
		log.Printf("Failed opening jpg file with path %s, error %s", fsPath, jpgErr)
	}
	defer jpgFile.Close()

	tgaFile, tgaErr := os.Open(tga)

	if tgaErr == nil {
		return true, fmt.Sprintf("%s.tga", material)
	} else {
		log.Printf("Failed opening tga file with path %s, error %s", fsPath, tgaErr)
	}

	defer tgaFile.Close()
	return false, material
}

func addTrailingSlash(path string) string {
	if path == "" {
		return path
	} else {
		rune := []rune(path)
		lastCharacter := string(rune[len(rune)-1])
		if strings.Contains(lastCharacter, "/") {
			return path
		} else {
			return path + "/"
		}
	}
}

func sortMaterials(materials map[string]int, basePath string) Materials {
	sorted := Materials{make(map[string]int), make(map[string]int)}
	for material := range materials {
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

func addTextureFileExtension(textures map[string]int, basePath string) map[string]int {
	returnValue := map[string]int{}
	for material := range textures {
		isT, filePath := isTexture(material, basePath)
		if isT {
			returnValue[filePath] = returnValue[filePath] + 1
		}
	}
	return returnValue
}
