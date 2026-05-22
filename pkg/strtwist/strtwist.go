package strtwist

import (
	"fmt"
	"strconv"
	"strings"
)

// Reverse returns the input string with its characters in reversed order.
func Reverse(input string) string {
	r := []rune(input)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
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

func k8sTransform(s string) string {
	if len(s) == 0 {
		return s
	}
	var first = string(s[0])
	var length = len(s)
	var last = string(s[length-1])
	if len(s) > 2 {
		return first + fmt.Sprint(length-2) + last
	}
	return s
}

// K8s returns the input string with each word shortened using Kubernetes-style numeronym notation.
func K8s(s string) string {
	list := strings.Split(s, " ")
	for i := 0; i < len(list); i++ {
		list[i] = k8sTransform(list[i])
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
