package tools

import "strings"

func Split(str string, step string) []string {
	if str == "" {
		return make([]string, 0)
	}
	return strings.Split(str, step)
}
