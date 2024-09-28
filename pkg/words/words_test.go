package words

import (
	"regexp"
	"testing"
)

func WordCount(value string) int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return len(results)
}

func TestWords(t *testing.T) {
	cases := []struct {
		in     int
		want   int
		weight WordSetWeight
	}{
		{24, 24, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}},
		{0, 0, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}},
		{1, 1, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}},
		{100, 0, WordSetWeight{}},
		{100, 100, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}},
		{100, 100, WordSetWeight{Animals: 1, Colors: 1}},
		{100, 100, WordSetWeight{Animals: 1}},
		{100, 100, WordSetWeight{Adjectives: 1, Animals: 10, Colors: 10, Nouns: 1, Verbs: 1}},
		{100, 100, WordSetWeight{Adjectives: 5, Animals: 4, Colors: 3, Nouns: 2, Verbs: 1}},
		{100, 100, WordSetWeight{Adjectives: 1, Animals: 0, Colors: 0, Nouns: 1, Verbs: 1}},
		{1000, 1000, WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}},
	}
	for _, c := range cases {
		got := WordCount(Words(c.in, c.weight))
		if got != c.want {
			t.Errorf("len(Words(%q)) == %q, want %q", c.in, got, c.want)
		}
	}
}
