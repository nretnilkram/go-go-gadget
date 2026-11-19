package go_go_gadget

import (
	"bytes"
	"strings"
	"testing"
	"unicode"
)

func execute(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	root := newRootCmd()
	var outBuf, errBuf bytes.Buffer
	root.SetOut(&outBuf)
	root.SetErr(&errBuf)
	root.SetArgs(args)
	err = root.Execute()
	return outBuf.String(), errBuf.String(), err
}

func TestNewRootCmd_RegistersExpectedCommands(t *testing.T) {
	root := newRootCmd()
	want := []string{
		"inspect", "k8s", "reverse", "symsub", "now", "words", "password",
		"m8s", "utilities", "grit", "group", "generate-documentation",
	}
	got := map[string]bool{}
	for _, cmd := range root.Commands() {
		got[cmd.Name()] = true
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("missing command %q", name)
		}
	}
}

func TestRoot_Version(t *testing.T) {
	out, _, err := execute(t, "--version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, version) {
		t.Fatalf("version output %q does not contain %q", out, version)
	}
}

func TestInspect(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "chars", args: []string{"inspect", "ab"}, want: "'ab' has 2 chars."},
		{name: "digits", args: []string{"inspect", "-d", "a1b2"}, want: "'a1b2' has 2 digits."},
		{name: "singular", args: []string{"inspect", "x"}, want: "'x' has 1 char."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, _, err := execute(t, tt.args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(out, tt.want) {
				t.Fatalf("output %q does not contain %q", out, tt.want)
			}
		})
	}
}

func TestK8s(t *testing.T) {
	out, _, err := execute(t, "k8s", "Internationalization")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(out); got != "I18n" {
		t.Fatalf("got %q, want I18n", got)
	}
}

func TestReverse(t *testing.T) {
	out, _, err := execute(t, "reverse", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(out); got != "cba" {
		t.Fatalf("got %q, want cba", got)
	}
}

func TestSymSub(t *testing.T) {
	out, _, err := execute(t, "symsub", "password")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := strings.TrimSpace(out); got != "p@$$w0rd" {
		t.Fatalf("got %q, want p@$$w0rd", got)
	}
}

func TestNow_Formats(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		check func(string) bool
	}{
		{
			name: "dash",
			args: []string{"now"},
			check: func(s string) bool {
				return strings.Contains(s, "-") && !strings.Contains(s, ":")
			},
		},
		{
			name: "unix",
			args: []string{"now", "--unix"},
			check: func(s string) bool {
				s = strings.TrimSpace(s)
				if s == "" {
					return false
				}
				for _, r := range s {
					if !unicode.IsDigit(r) {
						return false
					}
				}
				return true
			},
		},
		{
			name:  "colon",
			args:  []string{"now", "--colon"},
			check: func(s string) bool { return strings.Contains(s, ":") },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, _, err := execute(t, tt.args...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.check(out) {
				t.Fatalf("unexpected output: %q", out)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	out, _, err := execute(t, "password", "--length", "16", "--symbols=false")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimSpace(out)
	if len(got) != 16 {
		t.Fatalf("len = %d, want 16 (%q)", len(got), got)
	}
	for _, r := range got {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			t.Fatalf("unexpected rune %q in %q", r, got)
		}
	}
}

func TestWords(t *testing.T) {
	out, _, err := execute(t, "words", "--count", "3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) != 3 {
		t.Fatalf("got %d words, want 3: %q", len(fields), out)
	}
}

func TestSemver(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1.2.3", "true"},
		{"not-a-version", "false"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out, _, err := execute(t, "utilities", "semver", tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := strings.TrimSpace(out); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBumpChart_InvalidArgs(t *testing.T) {
	_, stderr, err := execute(t, "utilities", "bump-chart", "nope")
	if err == nil {
		t.Fatal("expected error for invalid bump type")
	}
	if !strings.Contains(stderr, "invalid bump type") {
		t.Fatalf("stderr %q missing invalid bump type message", stderr)
	}
}

func TestResolveImage(t *testing.T) {
	tests := []struct {
		name       string
		opts       m8sImageOptions
		wantImage  string
		wantPkgMgr string
	}{
		{name: "default alpine", opts: m8sImageOptions{}, wantImage: "alpine", wantPkgMgr: "apk"},
		{name: "busybox", opts: m8sImageOptions{useBusyBox: true}, wantImage: "busybox", wantPkgMgr: "none"},
		{name: "ubuntu", opts: m8sImageOptions{useUbuntu: true}, wantImage: "ubuntu", wantPkgMgr: "apt"},
		{name: "custom image wins", opts: m8sImageOptions{useUbuntu: true, containerImage: "debian"}, wantImage: "debian", wantPkgMgr: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			image, pkg := resolveImage(tt.opts)
			if image != tt.wantImage || pkg != tt.wantPkgMgr {
				t.Fatalf("resolveImage() = (%q, %q), want (%q, %q)", image, pkg, tt.wantImage, tt.wantPkgMgr)
			}
		})
	}
}

func TestInspect_RequiresArg(t *testing.T) {
	_, _, err := execute(t, "inspect")
	if err == nil {
		t.Fatal("expected error when inspect has no args")
	}
}

func TestGrit_MissingDir(t *testing.T) {
	t.Chdir(t.TempDir())
	_, _, err := execute(t, "grit", "status")
	if err == nil {
		t.Fatal("expected error outside a grit directory")
	}
}

func TestGroup_MissingDir(t *testing.T) {
	t.Chdir(t.TempDir())
	_, _, err := execute(t, "group", "ls")
	if err == nil {
		t.Fatal("expected error outside a grit directory")
	}
}

func TestDocs_ValidArgs(t *testing.T) {
	root := newRootCmd()
	docs, _, err := root.Find([]string{"docs"})
	if err != nil {
		t.Fatalf("find docs: %v", err)
	}
	want := map[string]bool{"md": true, "rest": true, "yaml": true}
	for _, arg := range docs.ValidArgs {
		if !want[arg] {
			t.Fatalf("unexpected ValidArg %q", arg)
		}
		delete(want, arg)
	}
	if len(want) != 0 {
		t.Fatalf("missing ValidArgs: %v", want)
	}
}
