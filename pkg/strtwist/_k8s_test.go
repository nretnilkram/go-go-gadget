package strtwist

import (
	"testing"
)

func TestK8s(t *testing.T) {
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
		got := K8s(c.in)
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
