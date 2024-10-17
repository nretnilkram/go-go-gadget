package now

import (
	"strings"
	"testing"
)

func TestShowDateTime(t *testing.T) {
	cases := []struct {
		in          string
		includeTime bool
		want        string
		minLength   int
		maxLength   int
	}{
		{"colon", false, ":", 8, 10},
		{"dash", false, "-", 8, 10},
		{"dash", true, "-", 14, 16},
		{"dot", false, ".", 8, 10},
		{"dot", true, ".", 14, 16},
		{"slash", false, "/", 8, 10},
		{"slash", true, "/", 14, 16},
	}
	for _, c := range cases {
		got := ShowDateTime(c.in, c.includeTime)
		if !strings.Contains(got, c.want) || !(len(got) >= c.minLength && len(got) <= c.maxLength) {
			t.Errorf("ShowDateTime(%q, %t) == %q, want %q", c.in, c.includeTime, got, c.want)
		}
	}
}
