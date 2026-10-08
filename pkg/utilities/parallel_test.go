package utilities

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommandParallel_Completes(t *testing.T) {
	tmp := t.TempDir()
	dirA := filepath.Join(tmp, "repo-a")
	dirB := filepath.Join(tmp, "repo-b")
	if err := os.Mkdir(dirA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(dirB, 0755); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		RunCommandParallel("pwd", nil, []ParallelDir{
			{Name: "repo-a", Path: dirA},
			{Name: "repo-b", Path: dirB},
		}, 0)
	})

	if !strings.Contains(out, "REPO-A") {
		t.Fatalf("expected REPO-A header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "REPO-B") {
		t.Fatalf("expected REPO-B header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[pwd]") {
		t.Fatalf("expected [pwd] command display in output, got:\n%s", out)
	}
}

func TestRunCommandParallel_WithArgs(t *testing.T) {
	tmp := t.TempDir()
	dirA := filepath.Join(tmp, "repo-a")
	if err := os.Mkdir(dirA, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dirA, "hello.txt"), []byte("hi"), 0644); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		RunCommandParallel("ls", []string{"hello.txt"}, []ParallelDir{
			{Name: "repo-a", Path: dirA},
		}, 0)
	})

	if !strings.Contains(out, "hello.txt") {
		t.Fatalf("expected hello.txt in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[ls hello.txt]") {
		t.Fatalf("expected [ls hello.txt] command display, got:\n%s", out)
	}
}

func TestRunCommandParallel_ReportsFailedDirs(t *testing.T) {
	tmp := t.TempDir()
	dirOK := filepath.Join(tmp, "repo-ok")
	if err := os.Mkdir(dirOK, 0755); err != nil {
		t.Fatal(err)
	}

	out := captureStdout(t, func() {
		RunCommandParallel("pwd", nil, []ParallelDir{
			{Name: "repo-ok", Path: dirOK},
			{Name: "repo-missing", Path: filepath.Join(tmp, "does-not-exist")},
		}, 0)
	})

	if !strings.Contains(out, "DIRECTORIES WITH ERRORS") {
		t.Fatalf("expected error summary in output, got:\n%s", out)
	}
	if !strings.Contains(out, "repo-missing:") {
		t.Fatalf("expected repo-missing in error summary, got:\n%s", out)
	}
	if strings.Contains(out, "  repo-ok:") {
		t.Fatalf("repo-ok should not appear in error summary, got:\n%s", out)
	}
}

func TestRunCommandParallel_SemaphoreLimitOne(t *testing.T) {
	tmp := t.TempDir()
	var dirs []ParallelDir
	for _, name := range []string{"a", "b", "c"} {
		path := filepath.Join(tmp, name)
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
		dirs = append(dirs, ParallelDir{Name: name, Path: path})
	}

	// Should complete without deadlock even with concurrency capped at 1.
	RunCommandParallel("pwd", nil, dirs, 1)
}

func TestRunCommandParallel_EmptyDirs(t *testing.T) {
	RunCommandParallel("pwd", nil, nil, 0)
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	var buf strings.Builder
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
