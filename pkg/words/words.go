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

	total := weight.Adjectives + weight.Animals + weight.Colors + weight.Nouns + weight.Verbs
	if total == 0 {
		return ""
	}

	wordSet := LoadJsonWords()

	// Cumulative upper-bound for each category. A random value n in [0, total)
	// maps to the first category whose threshold exceeds n, giving O(1) selection
	// instead of the O(totalWeight) slice that the previous implementation built
	// on every call.
	thresholds := [5]int{
		weight.Adjectives,
		weight.Adjectives + weight.Animals,
		weight.Adjectives + weight.Animals + weight.Colors,
		weight.Adjectives + weight.Animals + weight.Colors + weight.Nouns,
		total,
	}
	categories := [5][]string{
		wordSet.Adjectives,
		wordSet.Animals,
		wordSet.Colors,
		wordSet.Nouns,
		wordSet.Verbs,
	}

	var sb strings.Builder
	for i := 0; i < length; i++ {
		n := rand.Intn(total)
		for cat, threshold := range thresholds {
			if n < threshold {
				sb.WriteString(randomItem(categories[cat]))
				sb.WriteByte(' ')
				break
			}
		}
	}

	return sb.String()
}
