package util

import "strings"

func GetNameFromEmailAddress(address string) string {
	before, _, _ := strings.Cut(address, "@")
	return before
}
