package words

import (
	_ "embed"
	"encoding/json"
	"math/rand"
)

// WordSet holds categorized lists of English words loaded from the embedded JSON file.
type WordSet struct {
	Adjectives []string `json:"adjectives"`
	Animals    []string `json:"animals"`
	Colors     []string `json:"colors"`
	Nouns      []string `json:"nouns"`
	Verbs      []string `json:"verbs"`
}

// WordSetWeight controls how many words from each category are included in the weighted selection pool.
type WordSetWeight struct {
	Adjectives int
	Animals    int
	Colors     int
	Nouns      int
	Verbs      int
}

//go:embed english_words.json
var englishWords []byte

// LoadJsonWords unmarshals the embedded English word list JSON and returns it as a WordSet.
func LoadJsonWords() WordSet {
	var wordSet WordSet
	err := json.Unmarshal([]byte(englishWords), &wordSet)
	if err != nil {
		panic(err)
	}
	return wordSet
}

func randomItem(list []string) string {
	return list[rand.Intn(len(list))]
}

// Words returns a space-separated string of randomly selected words. The number of words
// is controlled by length, and the category distribution is determined by weight.
func Words(length int, weight WordSetWeight) string {
	if length == 0 {
		return ""
	}

	englishWords := LoadJsonWords()

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
