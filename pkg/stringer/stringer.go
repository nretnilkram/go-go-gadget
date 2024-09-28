package stringer

import (
	"fmt"
	"strconv"
	"strings"
)

// Revere string
func Reverse(input string) (result string) {
	for _, c := range input {
		result = string(c) + result
	}
	return result
}

func inspectNumbers(input string) (count int) {
	for _, c := range input {
		_, err := strconv.Atoi(string(c))
		if err == nil {
			count++
		}
	}
	return count
}

// Inspect a string and return character count
func Inspect(input string, digits bool) (count int, kind string) {
	if !digits {
		return len(input), "char"
	}
	return inspectNumbers(input), "digit"
}

func k8s_transform(s string) string {
	if len(s) == 0 {
		return s
	}
	var first string = string(s[0])
	var length int = len(s)
	var last string = string(s[length-1])
	if len(s) > 2 {
		return first + fmt.Sprint(length-2) + last
	}
	return s
}

// Kubernetes returns input string in the Kubernetes k8s shortened version.
func K8s(s string) string {
	list := strings.Split(s, " ")
	for i := 0; i < len(list); i++ {
		list[i] = k8s_transform(list[i])
	}
	return strings.Join(list, " ")
}

func subRune(r rune) rune {
	switch r {
	case 'a', 'A':
		return '@'
	case 'e', 'E':
		return '3'
	case 'i', 'I':
		return '!'
	case 'o', 'O':
		return '0'
	case 's', 'S':
		return '$'
	default:
		return r
	}
}

// SymbolSubstitution returns input string with certain values substituted.
func SymbolSubstitution(s string) string {
	return strings.Map(subRune, s)
}
