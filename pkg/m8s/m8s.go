package m8s

import (
	"strings"
)

func ReverseStrArray(target []string) []string {
	for front, back := 0, len(target)-1; front < back; front, back = front+1, back-1 {
		target[front], target[back] = target[back], target[front]
	}
	return target
}

func BreakOnSection(target string, separator string) string {
	parts := strings.Split(target, separator)
	final := ""

	for _, part := range ReverseStrArray(parts) {
		if (len(final) + len(part)) > 40 {
			break
		}
		final = part + separator + final
	}

	if final == "" {
		final = "long-image-name"
	}

	return strings.TrimRight(final, separator)
}

// Replace punctuation and truncate image to be used in the pod name
func Image2Name(imageName string, separator string) string {
	finalName := imageName

	finalName = strings.Replace(finalName, ".", separator, -1)
	finalName = strings.Replace(finalName, "/", separator, -1)
	finalName = strings.Replace(finalName, ":", separator, -1)
	finalName = strings.Replace(finalName, "-", separator, -1)

	if len(finalName) > 40 {
		finalName = BreakOnSection(finalName, separator)
	}

	return finalName
}
