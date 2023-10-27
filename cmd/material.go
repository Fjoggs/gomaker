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

func isTexture(material string, rootPath string) (bool, string) {
	path := addTrailingSlash(rootPath) + material
	jpg := fmt.Sprintf("%s.jpg", path)
	tga := fmt.Sprintf("%s.tga", path)
	jpgFile, jpgErr := os.Open(jpg)

	if jpgErr == nil {
		return true, jpg
	} else {
		log.Printf("Failed opening jpg file with path %s, error %s", path, jpgErr)
	}
	defer jpgFile.Close()

	tgaFile, tgaErr := os.Open(tga)

	if tgaErr == nil {
		return true, tga
	} else {
		log.Printf("Failed opening tga file with path %s, error %s", path, tgaErr)
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
