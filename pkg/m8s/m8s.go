package m8s

import (
	"strings"
)

func ReverseStrArray(arr []string) []string {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
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
	return strings.TrimRight(final, separator)
}

// Replace punctuation and truncate image to be used in the pod name
func Image2Name(imageName string) string {
	finalName := imageName

	finalName = strings.Replace(finalName, ".", "-", -1)
	finalName = strings.Replace(finalName, "/", "-", -1)
	finalName = strings.Replace(finalName, ":", "-", -1)

	if len(finalName) > 40 {
		finalName = BreakOnSection(finalName, "-")
	}

	return finalName
}
