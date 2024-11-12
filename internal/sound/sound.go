package builder

import (
	"fmt"
	"regexp"
	"strings"
)

func getSound(line string) string {
	soundFile := ""
	soundRegex := regexp.MustCompile(`((\w+[\/_-]*)+\/((\w)+[\/_-]*)*)+`)
	sound := soundRegex.FindString(line)
	if strings.Contains(line, ".wav") && !strings.Contains(line, "sound/world/") {
		soundFile = fmt.Sprintf("%s.wav", sound)
	}
	return soundFile
}
