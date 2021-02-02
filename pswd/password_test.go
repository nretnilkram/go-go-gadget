package pswd

import (
        "testing"
)

func TestPassword(t *testing.T) {
  cases := []struct {
    in int; want int; weight PasswordWeight
  }{
    {24, 24, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {3, 3, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {200, 200, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {2000, 2000, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {50, 0, PasswordWeight { lower: 0, upper: 0, digit: 0, symbol: 0, }},
    {10, 10, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {10, 10, PasswordWeight { lower: 1, upper: 1, digit: 1, symbol: 1, }},
    {10, 10, PasswordWeight { lower: 0, upper: 0, digit: 3, symbol: 3, }},
    {10, 10, PasswordWeight { lower: 1, upper: 0, digit: 0, symbol: 0, }},
    {100, 100, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
    {100, 100, PasswordWeight { lower: 100, upper: 100, digit: 100, symbol: 100, }},
    {0, 0, PasswordWeight { lower: 4, upper: 3, digit: 3, symbol: 2, }},
  }
  for _, c := range cases {
    got := len(Password(c.in, c.weight))
    if got != c.want {
      t.Errorf("len(Password(%q)) == %q, want %q", c.in, got, c.want)
    }
  }
}
