package password

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerate_Length(t *testing.T) {
	got, err := Generate(24, Weight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(got) != 24 {
		t.Fatalf("len = %d, want 24", len(got))
	}
}

func TestGenerate_NoSymbols(t *testing.T) {
	got, err := Generate(64, Weight{Lower: 4, Upper: 3, Digit: 3, Symbol: 0})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, r := range got {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			t.Fatalf("unexpected non-alphanumeric rune %q in %q", r, got)
		}
	}
}

func TestGenerate_WithSymbols(t *testing.T) {
	const alphabet = lowerChars + upperChars + digitChars + symbolChars
	got, err := Generate(48, Weight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	for _, r := range got {
		if !strings.ContainsRune(alphabet, r) {
			t.Fatalf("unexpected rune %q in %q", r, got)
		}
	}
}

func TestGenerate_InvalidLength(t *testing.T) {
	if _, err := Generate(0, Weight{Lower: 1}); err == nil {
		t.Fatal("expected error for length 0")
	}
}

func TestGenerate_NoWeights(t *testing.T) {
	if _, err := Generate(8, Weight{}); err == nil {
		t.Fatal("expected error when all weights are zero")
	}
}

func TestGenerate_NotDeterministic(t *testing.T) {
	a, err := Generate(32, Weight{Lower: 1, Upper: 1, Digit: 1, Symbol: 1})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate(32, Weight{Lower: 1, Upper: 1, Digit: 1, Symbol: 1})
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatalf("expected distinct passwords, both were %q", a)
	}
}
