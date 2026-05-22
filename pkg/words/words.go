package words

import (
	_ "embed"
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"sync"
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
var embeddedWords []byte

var (
	loadOnce    sync.Once
	cachedWords WordSet
)

// LoadJsonWords unmarshals the embedded English word list JSON and returns it as a WordSet.
// The JSON is parsed only once; subsequent calls return the cached result.
func LoadJsonWords() WordSet {
	loadOnce.Do(func() {
		if err := json.Unmarshal(embeddedWords, &cachedWords); err != nil {
			log.Fatalf("failed to parse word list: %v", err)
		}
	})
	return cachedWords
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

	wordSet := LoadJsonWords()

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

	var sb strings.Builder
	for i := 0; i < length; i++ {
		switch weighted[rand.Intn(len(weighted))] {
		case 0:
			sb.WriteString(randomItem(wordSet.Adjectives) + " ")
		case 1:
			sb.WriteString(randomItem(wordSet.Animals) + " ")
		case 2:
			sb.WriteString(randomItem(wordSet.Colors) + " ")
		case 3:
			sb.WriteString(randomItem(wordSet.Nouns) + " ")
		case 4:
			sb.WriteString(randomItem(wordSet.Verbs) + " ")
		}
	}

	return sb.String()
}
