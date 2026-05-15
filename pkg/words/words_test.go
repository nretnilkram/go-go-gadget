package words

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"
)

func WordCount(value string) int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return len(results)
}

func TestLoadJsonWords(t *testing.T) {
	mockData := `{
		"adjectives": ["quick", "lazy"],
		"animals": ["fox", "dog"],
		"colors": ["brown", "red"],
		"nouns": ["jumps", "runs"],
		"verbs": ["over", "under"]
	}`

	var wordSet WordSet
	err := json.Unmarshal([]byte(mockData), &wordSet)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock data: %v", err)
	}

	if len(wordSet.Adjectives) != 2 || wordSet.Adjectives[0] != "quick" {
		t.Errorf("Expected adjectives to contain 'quick', got %v", wordSet.Adjectives)
	}
	if len(wordSet.Animals) != 2 || wordSet.Animals[0] != "fox" {
		t.Errorf("Expected animals to contain 'fox', got %v", wordSet.Animals)
	}
	if len(wordSet.Colors) != 2 || wordSet.Colors[0] != "brown" {
		t.Errorf("Expected colors to contain 'brown', got %v", wordSet.Colors)
	}
	if len(wordSet.Nouns) != 2 || wordSet.Nouns[0] != "jumps" {
		t.Errorf("Expected nouns to contain 'jumps', got %v", wordSet.Nouns)
	}
	if len(wordSet.Verbs) != 2 || wordSet.Verbs[0] != "over" {
		t.Errorf("Expected verbs to contain 'over', got %v", wordSet.Verbs)
	}
}

func TestRandomItem(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	for i := 0; i < 10; i++ {
		item := randomItem(items)
		if item != "apple" && item != "banana" && item != "cherry" {
			t.Errorf("Expected item to be one of 'apple', 'banana', or 'cherry', got %v", item)
		}
	}
}

func TestWords(t *testing.T) {
	mockData := `{
		"adjectives": ["quick", "lazy"],
		"animals": ["fox", "dog"],
		"colors": ["brown", "red"],
		"nouns": ["jumps", "runs"],
		"verbs": ["over", "under"]
	}`

	var wordSet WordSet
	err := json.Unmarshal([]byte(mockData), &wordSet)
	if err != nil {
		t.Fatalf("Failed to unmarshal mock data: %v", err)
	}

	cases := []struct {
		length int
		weight WordSetWeight
		want   int
	}{
		{5, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 5},
		{0, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 0},
		{10, WordSetWeight{Adjectives: 2, Animals: 2, Colors: 2, Nouns: 2, Verbs: 2}, 10},
		{3, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1}, 3},
		{24, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 24},
		{0, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 0},
		{1, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 1},
		{100, WordSetWeight{}, 0},
		{100, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 100},
		{100, WordSetWeight{Animals: 1, Colors: 1}, 100},
		{100, WordSetWeight{Animals: 1}, 100},
		{100, WordSetWeight{Adjectives: 1, Animals: 10, Colors: 10, Nouns: 1, Verbs: 1}, 100},
		{100, WordSetWeight{Adjectives: 5, Animals: 4, Colors: 3, Nouns: 2, Verbs: 1}, 100},
		{100, WordSetWeight{Adjectives: 1, Animals: 0, Colors: 0, Nouns: 1, Verbs: 1}, 100},
		{1000, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}, 1000},
	}

	for _, c := range cases {
		got := WordCount(Words(c.length, c.weight))
		if got != c.want {
			t.Errorf("WordCount(Words(%d, %v)) == %d, want %d", c.length, c.weight, got, c.want)
		}
	}
}

// TestWordsWeightAllZeroReturnsEmpty confirms the zero-weight boundary: when every
// category weight is 0 the weighted pool is empty and Words must return "".
func TestWordsWeightAllZeroReturnsEmpty(t *testing.T) {
	result := Words(100, WordSetWeight{})
	if result != "" {
		t.Errorf("Words(100, WordSetWeight{}) = %q, want empty string", result)
	}
}

// TestWordsWeightSingleCategory confirms that when exactly one category carries a
// non-zero weight every word in the output belongs to that category.
func TestWordsWeightSingleCategory(t *testing.T) {
	ws := LoadJsonWords()

	cases := []struct {
		name     string
		weight   WordSetWeight
		validSet []string
	}{
		{"adjectives only", WordSetWeight{Adjectives: 1}, ws.Adjectives},
		{"animals only", WordSetWeight{Animals: 1}, ws.Animals},
		{"colors only", WordSetWeight{Colors: 1}, ws.Colors},
		{"nouns only", WordSetWeight{Nouns: 1}, ws.Nouns},
		{"verbs only", WordSetWeight{Verbs: 1}, ws.Verbs},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			valid := make(map[string]bool, len(c.validSet))
			for _, w := range c.validSet {
				valid[w] = true
			}

			for _, word := range strings.Fields(Words(50, c.weight)) {
				if !valid[word] {
					t.Errorf("word %q is not in the expected category", word)
				}
			}
		})
	}
}

// TestWordsWeightExcludesZeroWeightCategories confirms that categories with a zero
// weight never contribute words to the output, even when other categories are active.
// To avoid false positives from words that appear in multiple categories, only words
// that are exclusively in excluded categories (not present in any included category)
// are checked.
func TestWordsWeightExcludesZeroWeightCategories(t *testing.T) {
	ws := LoadJsonWords()

	cases := []struct {
		name     string
		weight   WordSetWeight
		included [][]string
		excluded [][]string
	}{
		{
			name:     "only adjectives and nouns active",
			weight:   WordSetWeight{Adjectives: 1, Nouns: 1},
			included: [][]string{ws.Adjectives, ws.Nouns},
			excluded: [][]string{ws.Animals, ws.Colors, ws.Verbs},
		},
		{
			name:     "only animals active",
			weight:   WordSetWeight{Animals: 5},
			included: [][]string{ws.Animals},
			excluded: [][]string{ws.Adjectives, ws.Colors, ws.Nouns, ws.Verbs},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Build a set of words present in at least one included category.
			includedSet := make(map[string]bool)
			for _, list := range c.included {
				for _, w := range list {
					includedSet[w] = true
				}
			}

			// Words that appear only in excluded categories (never in included ones).
			exclusivelyExcluded := make(map[string]bool)
			for _, list := range c.excluded {
				for _, w := range list {
					if !includedSet[w] {
						exclusivelyExcluded[w] = true
					}
				}
			}

			for _, word := range strings.Fields(Words(200, c.weight)) {
				if exclusivelyExcluded[word] {
					t.Errorf("word %q from a zero-weight category appeared in output", word)
				}
			}
		})
	}
}
