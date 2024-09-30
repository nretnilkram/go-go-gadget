package strtwist

import (
	"testing"
)

func TestStringInspect(t *testing.T) {
	cases := []struct {
		inStr       string
		inDigitBool bool
		wantCount   int
		wantString  string
	}{
		{"abcdefghijklmnopqrstuvwxyz", false, 26, "char"},
		{"Mark", false, 4, "char"},
		{"Nretnil Kram", false, 12, "char"},
		{"a1b2c3d4 e5f6 g7", false, 16, "char"},
		{"abcdefghijklmnopqrstuvwxyz", true, 0, "digit"},
		{"abcdefghijklmnopqrstuvwxyz0123456789", true, 10, "digit"},
		{"a1b2c3d4 e5f6 g7", true, 7, "digit"},
	}
	for _, c := range cases {
		gotCount, gotString := Inspect(c.inStr, c.inDigitBool)
		if gotCount != c.wantCount || gotString != c.wantString {
			t.Errorf("Inspect(%q, %t) == (%q, %q), want (%q, %q)", c.inStr, c.inDigitBool, c.wantString, c.wantCount, gotCount, gotString)
		}
	}
}
