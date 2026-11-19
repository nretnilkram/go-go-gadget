package m8s

import (
	"fmt"
	"io"
	"strings"
)

// ReverseStrArray reverses the given string slice in place and returns it.
func ReverseStrArray(target []string) []string {
	for front, back := 0, len(target)-1; front < back; front, back = front+1, back-1 {
		target[front], target[back] = target[back], target[front]
	}
	return target
}

// PrintInfo writes a formatted summary of the Kubernetes resource details to w.
func PrintInfo(w io.Writer, k8sType string, imageName string, pkgManager string, resourceName string) {
	_, _ = fmt.Fprintln(w, "\n##############")
	_, _ = fmt.Fprintf(w, "# Type: %s\n", k8sType)
	_, _ = fmt.Fprintf(w, "# Image: %s\n", imageName)
	_, _ = fmt.Fprintf(w, "# Pkg Manager: %s\n", pkgManager)
	_, _ = fmt.Fprintf(w, "# Name: %s\n", resourceName)
	_, _ = fmt.Fprintf(w, "##############\n\n")
}

// BreakOnSection truncates target by splitting on separator and keeping only the trailing
// sections that fit within 40 characters.
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

// Image2Name sanitizes a container image name for use as a Kubernetes pod name by replacing
// punctuation with the given separator and truncating to 40 characters.
func Image2Name(imageName string, separator string) string {
	finalName := imageName

	finalName = strings.ReplaceAll(finalName, ".", separator)
	finalName = strings.ReplaceAll(finalName, "/", separator)
	finalName = strings.ReplaceAll(finalName, ":", separator)
	finalName = strings.ReplaceAll(finalName, "-", separator)

	if len(finalName) > 40 {
		finalName = BreakOnSection(finalName, separator)
	}

	return finalName
}
