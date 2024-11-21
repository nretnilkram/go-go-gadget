package m8s

import (
	"testing"
)

func checkArrayEquality(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestM8sReverseStrArray(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{[]string{"a", "b", "c", "d"}, []string{"d", "c", "b", "a"}},
	}
	for _, c := range cases {
		got := ReverseStrArray(c.in)
		if !checkArrayEquality(got, c.want) {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestM8sBreakOnSection(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"555752902066-dkr-ecr-us-east-2-amazonaws-com-docker-hub-curlimages-curl-latest", "com-docker-hub-curlimages-curl-latest"},
		{"curlimages-curl-latest", "curlimages-curl-latest"},
	}
	for _, c := range cases {
		got := BreakOnSection(c.in, "-")
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestM8sImage2Name(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"555752902066.dkr.ecr.us-east-2.amazonaws.com/docker-hub/curlimages/curl:latest", "com-docker-hub-curlimages-curl-latest"},
		{"curlimages/curl:latest", "curlimages-curl-latest"},
	}
	for _, c := range cases {
		got := Image2Name(c.in)
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
