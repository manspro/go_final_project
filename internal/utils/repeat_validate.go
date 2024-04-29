package utils

import "strings"

func IsValidRepeatFormat(repeat string) bool {
	return strings.HasPrefix(repeat, "d ")
}
