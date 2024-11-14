package sound

import (
	"fmt"
	"regexp"
	"strings"
)

func GetSound(line string) string {
	soundFile := ""
	soundRegex := regexp.MustCompile(`((\w+[\/_-]*)+\/((\w)+[\/_-]*)*)+`)
	sound := soundRegex.FindString(line)
	if strings.Contains(line, ".wav") && !strings.Contains(line, "sound/world/") {
		soundFile = fmt.Sprintf("%s.wav", sound)
	}
	return soundFile
}

func AddSounds(line string, sounds map[string]int) {
	sound := GetSound(line)
	if len(sound) > 0 {
		sounds[sound] = sounds[sound] + 1
	}
}
