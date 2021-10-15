// Package strtwist implements additional functions UTF-8 encoded strings.
package strtwist

import (
	"strings"
)

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
