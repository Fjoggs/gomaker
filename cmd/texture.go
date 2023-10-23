package gomaker

import (
	"fmt"
	"regexp"
	"strings"
)

func isTexture(line string) string {
	textureRegex := regexp.MustCompile(`((\w+[\/_-]*)+\/((\w)+[\/_-]*)*)+`)
	texture := textureRegex.FindString(line)
	fmt.Println(texture)
	if len(texture) > 0 {
		if isCustomTexture(texture) {
			return formatPath(texture)
		}
	}
	return ""
}

func isCustomTexture(texture string) bool {
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
