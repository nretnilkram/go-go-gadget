package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

	cmdTest := RunCommand("ls", []string{"-lah"}, ".")
	if reflect.TypeOf(cmdTest).Kind() != reflect.String {
		t.Errorf("RunCommand is not a string.")
	}

	igr := IsGitRepo(".")
	if reflect.TypeOf(igr).Kind() != reflect.Bool {
		t.Errorf("IsGitRepo did not return a boolean.")
	}
	if !igr {
		t.Errorf("IsGitRepo did not return true for a git repo.")
	}

	igr2 := IsGitRepo(t.TempDir())
	if igr2 {
		t.Errorf("IsGitRepo did not return false for a non-git directory.")
	}
}

func TestFileDirExistsNonExistent(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	exists, err := FileDirExists(missing)
	if err != nil {
		t.Errorf("FileDirExists returned unexpected error: %v", err)
	}
	if exists {
		t.Errorf("FileDirExists returned true for a path that does not exist")
	}
}

func TestRunCommandError(t *testing.T) {
	result := RunCommand("/nonexistent/binary-xyz", []string{}, ".")
	if result == "" {
		t.Errorf("RunCommand should return a non-empty string on error")
	}
}

func TestShowDateTimeDefault(t *testing.T) {
	got := ShowDateTime("bogus", false)
	if got == "" {
		t.Errorf("ShowDateTime with unknown format returned empty string")
	}
}

func TestGrepFileForTFResources(t *testing.T) {
	content := `resource "aws_s3_bucket" "my_bucket" {
module "my_module" {
  source = "./modules/foo"
}
`
	tmp := filepath.Join(t.TempDir(), "main.tf")
	if err := os.WriteFile(tmp, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got := GrepFileForTFResources(tmp)

	if len(got) != 2 {
		t.Fatalf("expected 2 targets, got %d: %v", len(got), got)
	}

	cases := []string{
		"--target=aws_s3_bucket.my_bucket",
		"--target=module.my_module",
	}
	for i, want := range cases {
		if !strings.Contains(got[i], want) {
			t.Errorf("got[%d] = %q, want it to contain %q", i, got[i], want)
		}
	}
}

func TestListTFResources(t *testing.T) {
	content := fmt.Sprintf("resource \"aws_instance\" \"web\" {\n")
	tmp := filepath.Join(t.TempDir(), "main.tf")
	if err := os.WriteFile(tmp, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	ListTFResources([]string{tmp})
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
		{"dash", true, "-", 14, 19},
		{"dot", false, ".", 8, 10},
		{"dot", true, ".", 14, 19},
		{"slash", false, "/", 8, 10},
		{"slash", true, "/", 14, 19},
	}
	for _, c := range cases {
		got := ShowDateTime(c.in, c.includeTime)
		if !strings.Contains(got, c.want) || (len(got) < c.minLength || len(got) > c.maxLength) {
			t.Errorf("ShowDateTime(%q, %t) == %q, want %q", c.in, c.includeTime, got, c.want)
		}
	}
}
