package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestGoGoGadgetReverse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
		{"A", "A"},
		{"aBc", "cBa"},
		{"Mark", "kraM"},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "reverse", c.in)
		out, err := cmd.CombinedOutput()
		got := strings.TrimSuffix(string(out), "\n") // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget reverse %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGoGoGadgetSymSub(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"abcdefghijklmnopqrstuvwxyz", "@bcd3fgh!jklmn0pqr$tuvwxyz"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "@BCD3FGH!JKLMN0PQR$TUVWXYZ"},
		{"ABCDEF GHIJKLMNOP QRSTUVW XYZ", "@BCD3F GH!JKLMN0P QR$TUVW XYZ"},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "symsub", c.in)
		out, err := cmd.CombinedOutput()
		got := strings.TrimSuffix(string(out), "\n") // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget symsub %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGoGoGadgetInspect(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"abcdefghijklmnopqrstuvwxyz", "'abcdefghijklmnopqrstuvwxyz' has 26 chars."},
		{"Mark", "'Mark' has 4 chars."},
		{"Nretnil Kram", "'Nretnil Kram' has 12 chars."},
		{"a1b2c3d4 e5f6 g7", "'a1b2c3d4 e5f6 g7' has 16 chars."},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "inspect", c.in)
		out, err := cmd.CombinedOutput()
		got := strings.TrimSuffix(string(out), "\n") // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget inspect %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGoGoGadgetInspectDigits(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"abcdefghijklmnopqrstuvwxyz", "'abcdefghijklmnopqrstuvwxyz' has 0 digits."},
		{"abcdefghijklmnopqrstuvwxyz0123456789", "'abcdefghijklmnopqrstuvwxyz0123456789' has 10 digits."},
		{"a1b2c3d4 e5f6 g7", "'a1b2c3d4 e5f6 g7' has 7 digits."},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "inspect", "--digits", c.in)
		out, err := cmd.CombinedOutput()
		got := strings.TrimSuffix(string(out), "\n") // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget inspect %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGoGoGadgetK8s(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Kubernetes", "K8s"},
		{"kubernetes", "k8s"},
		{"kuberneteS", "k8S"},
		{"KUBERNETES", "K8S"},
		{"Mark", "M2k"},
		{"", ""},
		{"xy", "xy"},
		{"xyz", "x1z"},
		{"one two three", "o1e t1o t3e"},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "k8s", c.in)
		out, err := cmd.CombinedOutput()
		got := strings.TrimSuffix(string(out), "\n") // because out is []byte with a new line
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget k8s %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestGoGoGadgetPassword(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{8, 8},
		{10, 10},
		{1000, 1000},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "password", "--length", strconv.Itoa(c.in))
		out, err := cmd.CombinedOutput()
		got := len(strings.TrimSuffix(string(out), "\n")) // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget password %q == %q, want %q", c.in, got, c.want)
		}
	}
}

func WordCount(value string) int {
	// Match non-space character sequences.
	re := regexp.MustCompile(`[\S]+`)

	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return len(results)
}

func TestGoGoGadgetWords(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{8, 8},
		{10, 10},
		{1000, 1000},
	}
	for _, c := range cases {
		cmd := exec.Command("go-go-gadget", "words", "--count", strconv.Itoa(c.in))
		out, err := cmd.CombinedOutput()
		got := WordCount(strings.TrimSuffix(string(out), "\n")) // because out is []byte
		if err != nil || got != c.want {
			fmt.Println(got, out, err)
			t.Errorf("go-go-gadget symsub %q == %q, want %q", c.in, got, c.want)
		}
	}
}
