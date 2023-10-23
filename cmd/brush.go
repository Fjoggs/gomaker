package gomaker

import (
	"log"
	"regexp"
	"strings"
)

func isBrush(line string) bool {
	match, err := regexp.MatchString("// brush", strings.ToLower(line))

	if err != nil {
		log.Fatal("Something went wrong", err)
	}

	return match
}
