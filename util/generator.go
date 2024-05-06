package util

import "strings"

func GenerateSlug(value string) string {
	return strings.ReplaceAll(value, " ", "-")
}