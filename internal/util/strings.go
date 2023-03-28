package util

import (
	"regexp"
	"strings"
)

var spaceEater = regexp.MustCompile(`\s+`)

func CompactString(q string) string {
	q = strings.TrimSpace(q)
	return spaceEater.ReplaceAllString(q, " ")
}
