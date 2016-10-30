package helpers

import (
	"strings"
)

func CountWords(word, content string) int {
	count := strings.Count(content, word)
	return count
}

func IsLogginPage(page string) bool {
	if strings.Contains(page, "loginForm") {
		return true
	}
	return false
}

// remove spaces, dots and carriage returns
func RemoveNoise(value string) string {
	temp := strings.Replace(value, " ", "", -1)
	temp = strings.Replace(temp, ".", "", -1)
	return strings.Replace(temp, "\n", "", -1)
}
