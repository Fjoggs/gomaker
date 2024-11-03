package gomaker

import (
	"fmt"
	"regexp"
	"strings"
)

func isBrush(line string) bool {
	match, err := regexp.MatchString("// brush", strings.ToLower(line))

	if err != nil {
		fmt.Println("Something went wrong", err)
	}

	return match
}
