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

func randSymbol() string {
  char_array := []rune(Symbols)
  rand.Seed(time.Now().UnixNano())
  return string(char_array[rand.Intn(len(char_array))])
}

func randDigit() string {
  char_array := []rune(Digits)
  rand.Seed(time.Now().UnixNano())
  return string(char_array[rand.Intn(len(char_array))])
}

func randLower() string {
  char_array := []rune(LowerLetters)
  rand.Seed(time.Now().UnixNano())
  return string(char_array[rand.Intn(len(char_array))])
}

func randUpper() string {
  char_array := []rune(UpperLetters)
  rand.Seed(time.Now().UnixNano())
  return string(char_array[rand.Intn(len(char_array))])
}

type PasswordWeight struct {
    lower int
    upper int
    digit int
    symbol int
}

// Password returns a password based on its input parameters.
func Password(length int, weight PasswordWeight) string {
  if length == 0 {
    return ""
  }

  var weighted []int
  for i := 0; i < weight.lower; i++ {
    weighted = append(weighted, 0)
  }
  for i := 0; i < weight.upper; i++ {
    weighted = append(weighted, 1)
  }
  for i := 0; i < weight.digit; i++ {
    weighted = append(weighted, 2)
  }
  for i := 0; i < weight.symbol; i++ {
    weighted = append(weighted, 3)
  }

  if len(weighted) == 0 {
    return ""
  }

  var password string
  for i := 0; i < length; i++ {
    rand.Seed(time.Now().UnixNano())
    switch weighted[rand.Intn(len(weighted))] {
    case 0:
      password += randLower()
    case 1:
      password += randUpper()
    case 2:
      password += randDigit()
    case 3:
      password += randSymbol()
    }
  }
  return password
}
