// Package pswd implements password generation functionality.
package words

import (
    "encoding/json"
    "fmt"
    "math/rand"
    "time"
    "os"
)

type WordSet struct {
    Adjectives  []string `json:"adjectives"`
    Animals     []string `json:"animals"`
    Colors      []string `json:"colors"`
    Nouns       []string `json:"nouns"`
    Verbs       []string `json:"verbs"`
}

type WordSetWeight struct {
    Adjectives  int
    Animals     int
    Colors      int
    Nouns       int
    Verbs       int
}

func LoadJsonWords (filename string) WordSet {
  var wordSet WordSet
  jsonFile, err := os.Open(filename)
  defer jsonFile.Close()
  if err != nil {
    fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(jsonFile)
  jsonParser.Decode(&wordSet)
  return wordSet
}

func randomItem (list []string) string {
  rand.Seed(time.Now().UnixNano())
  return list[rand.Intn(len(list))]
}

func getEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func Words(length int, weight WordSetWeight) string {
  if length == 0 {
    return ""
  }

  words_file := getEnv("GO_GO_GADGET_WORDS_FILE", "english_words.json")

  englishWords := LoadJsonWords(words_file)

  var weighted []int
  for i := 0; i < weight.Adjectives; i++ {
    weighted = append(weighted, 0)
  }
  for i := 0; i < weight.Animals; i++ {
    weighted = append(weighted, 1)
  }
  for i := 0; i < weight.Colors; i++ {
    weighted = append(weighted, 2)
  }
  for i := 0; i < weight.Nouns; i++ {
    weighted = append(weighted, 3)
  }
  for i := 0; i < weight.Verbs; i++ {
    weighted = append(weighted, 4)
  }

  if len(weighted) == 0 {
    return ""
  }

  var word_list string
  for i := 0; i < length; i++ {
    rand.Seed(time.Now().UnixNano())
    switch weighted[rand.Intn(len(weighted))] {
    case 0:
      word_list += randomItem(englishWords.Adjectives) + " "
    case 1:
      word_list += randomItem(englishWords.Animals) + " "
    case 2:
      word_list += randomItem(englishWords.Colors) + " "
    case 3:
      word_list += randomItem(englishWords.Nouns) + " "
    case 4:
      word_list += randomItem(englishWords.Verbs) + " "
    }
  }

  return word_list
}
