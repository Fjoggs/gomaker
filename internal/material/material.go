package material

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Materials struct {
	Textures map[string]int
	Shaders  map[string]int
}

func GetMaterial(line string) string {
	textureRegex := regexp.MustCompile(`((\w+[\/_-]*)+\/((\w)+[\/_-]*)*)+`)
	texture := textureRegex.FindString(line)
	if len(texture) > 0 {
		if IsCustomMaterial(texture) {
			return FormatPath(texture)
		}
	}
	return ""
}

func IsCustomMaterial(texture string) bool {
	common := [5]string{"common/", "common_alphascale/", "sfx/", "liquids/", "effects/"}
	for _, value := range common {
		if strings.Contains(texture, value) {
			return false
		}
	}
	return true
}

func FormatPath(texture string) string {
	return strings.Replace(texture, "textures/", "", 1)
}

func IsTexture(material string, baseFolderPath string) (bool, string) {
	fsPath := AddTrailingSlash(baseFolderPath) + "textures/" + material
	jpg := fmt.Sprintf("%s.jpg", fsPath)
	tga := fmt.Sprintf("%s.tga", fsPath)
	jpgFile, jpgErr := os.Open(jpg)

	if jpgErr == nil {
		return true, fmt.Sprintf("textures/%s.jpg", material)
	}

	defer jpgFile.Close()

	tgaFile, tgaErr := os.Open(tga)

	if tgaErr == nil {
		return true, fmt.Sprintf("textures/%s.tga", material)
	}

	defer tgaFile.Close()
	return false, material
}

func AddTrailingSlash(path string) string {
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

func SortMaterials(materials map[string]int, basePath string) Materials {
	sorted := Materials{make(map[string]int), make(map[string]int)}
	for material := range materials {
		isT, filePath := IsTexture(material, basePath)
		if isT {
			sorted.Textures[filePath] = sorted.Textures[filePath] + 1
			// It can also be a shader
			sorted.Shaders[material] = sorted.Shaders[material] + 1
		} else {
			sorted.Shaders[material] = sorted.Shaders[material] + 1
		}
	}
	return sorted
}

func AddTexturePathWithExtension(textures map[string]int, basePath string) map[string]int {
	returnValue := map[string]int{}
	for material := range textures {
		isT, filePath := IsTexture(material, basePath)
		if isT {
			returnValue[filePath] = returnValue[filePath] + 1
		}
	}
	return returnValue
}
