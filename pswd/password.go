// Package pswd implements password generation functionality.
package pswd

import (
        "math/rand"
        "time"
)

const (
  LowerLetters = "abcdefghijklmnopqrstuvwxyz"
  UpperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
  Digits = "0123456789"
  Symbols = "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"
)

// Password returns a password based on its input parameters.
func Password(length int, symbols bool) string {
  characters := LowerLetters + UpperLetters + Digits
  if symbols {
    characters += Symbols
  }
  char_array := []rune(characters)
  password := ""
  for i := 0; i < length; i++ {
    rand.Seed(time.Now().UnixNano())
    password += string(characters[rand.Intn(len(char_array))])
  }
  return password
}
