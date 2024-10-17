package utilities

import (
	"reflect"
	"testing"
)

func TestUtilities(t *testing.T) {
	gwd := GetWorkingDir()
	if reflect.TypeOf(gwd).Kind() != reflect.String {
		t.Errorf("GetWorkingDir did not return a string.")
	}

	doesIt, _ := FileDirExists(".")
	if reflect.TypeOf(doesIt).Kind() != reflect.Bool {
		t.Errorf("FileDirExists did not return a boolean.")
	}
	if !doesIt {
		t.Errorf("FileDirExists not true.")
	}

	cmdTest := RunShellCommand("ls -lah", ".")
	if reflect.TypeOf(cmdTest).Kind() != reflect.String {
		t.Errorf("RunShellCommand is not a string.")
	}

	igr := IsGitRepo(".")
	if reflect.TypeOf(igr).Kind() != reflect.Bool {
		t.Errorf("GetWorkingDir did not return a boolean.")
	}
	if !igr {
		t.Errorf("Did not return that this is a git repo.")
	}

	igr2 := IsGitRepo("../../..")
	if igr2 {
		t.Errorf("Did not return that parent directory is not a git repo.")
	}
}

func TestRegexTest(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"0.0.0", true},
		{"1", false},
		{"1.0", false},
		{"1.0.0", true},
		{"0.0.0-beta.1", true},
		{"0.0.0-go-go.RBTwed1HNVz6qtf59gO2WsdPVsRyGldAXxSv1P3ke0ogrjo6W71N5Zy5gzC8wy7J", true},
		{"2024.10.17", true},
	}
	for _, c := range cases {
		got := RegexTest(c.in, SemverRegex)
		if got != c.want {
			t.Errorf("K8s(%q) == %t, want %t", c.in, got, c.want)
		}
	}
}
