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
		{[]string{"one", "two", "three", "four"}, []string{"four", "three", "two", "one"}},
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
		in, want, separator string
	}{
		{"555752902066-dkr-ecr-us-east-2-amazonaws-com-docker-hub-curlimages-curl-latest", "com-docker-hub-curlimages-curl-latest", "-"},
		{"555752902066xdkrxecrxusxeastx2xamazonawsxcomxdockerxhubxcurlimagesxcurlxlatest", "comxdockerxhubxcurlimagesxcurlxlatest", "x"},
		{"curlimages-curl-latest", "curlimages-curl-latest", "-"},
		{"curlimagesxcurlxlatest", "curlimagesxcurlxlatest", "x"},
		{"curlimages-curl-verylongsupercalifragilisticexpialidocious", "long-image-name", "-"},
	}
	for _, c := range cases {
		got := BreakOnSection(c.in, c.separator)
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestM8sImage2Name(t *testing.T) {
	cases := []struct {
		in, want, separator string
	}{
		{"555752902066.dkr.ecr.us-east-2.amazonaws.com/docker-hub/curlimages/curl:latest", "com-docker-hub-curlimages-curl-latest", "-"},
		{"555752902066.dkr.ecr.us-east-2.amazonaws.com/docker-hub/curlimages/curl:latest", "comxdockerxhubxcurlimagesxcurlxlatest", "x"},
		{"555752902066.dkr.ecr.us-east-2.amazonaws.com/docker-hub/curlimages/curl:latest", "com_docker_hub_curlimages_curl_latest", "_"},
		{"curlimages/curl:latest", "curlimages-curl-latest", "-"},
		{"curlimages/curl:latest", "curlimagesxcurlxlatest", "x"},
		{"curlimages/curl:latest", "curlimages_curl_latest", "_"},
	}
	for _, c := range cases {
		got := Image2Name(c.in, c.separator)
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
