package pswd

import (
        "testing"
)

func TestPassword(t *testing.T) {
  cases := []struct {
    in, want int
  }{
    {24, 24},
    {10, 10},
    {100, 100},
  }
  for _, c := range cases {
    got := len(Password(c.in, true))
    if got != c.want {
      t.Errorf("len(Password(%q)) == %q, want %q", c.in, got, c.want)
    }
  }
}
