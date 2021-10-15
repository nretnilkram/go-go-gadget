package pswd

import (
	"testing"
)

func TestPassword(t *testing.T) {
	cases := []struct {
		in     int
		want   int
		weight PasswordWeight
	}{
		{24, 24, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{3, 3, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{200, 200, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{2000, 2000, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{50, 0, PasswordWeight{Lower: 0, Upper: 0, Digit: 0, Symbol: 0}},
		{10, 10, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{10, 10, PasswordWeight{Lower: 1, Upper: 1, Digit: 1, Symbol: 1}},
		{10, 10, PasswordWeight{Lower: 0, Upper: 0, Digit: 3, Symbol: 3}},
		{10, 10, PasswordWeight{Lower: 1, Upper: 0, Digit: 0, Symbol: 0}},
		{100, 100, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
		{100, 100, PasswordWeight{Lower: 100, Upper: 100, Digit: 100, Symbol: 100}},
		{0, 0, PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}},
	}
	for _, c := range cases {
		got := len(Password(c.in, c.weight))
		if got != c.want {
			t.Errorf("len(Password(%q)) == %q, want %q", c.in, got, c.want)
		}
	}
}
