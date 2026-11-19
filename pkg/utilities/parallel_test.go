package utilities

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultMaxConcurrent(t *testing.T) {
	got := DefaultMaxConcurrent()
	if got < 1 || got > maxConcurrentCap {
		t.Fatalf("DefaultMaxConcurrent() = %d, want in [1, %d]", got, maxConcurrentCap)
	}
}

func TestResolveMaxConcurrent(t *testing.T) {
	if got := resolveMaxConcurrent(0); got != DefaultMaxConcurrent() {
		t.Fatalf("resolveMaxConcurrent(0) = %d, want %d", got, DefaultMaxConcurrent())
	}
	if got := resolveMaxConcurrent(-1); got != DefaultMaxConcurrent() {
		t.Fatalf("resolveMaxConcurrent(-1) = %d, want %d", got, DefaultMaxConcurrent())
	}
	if got := resolveMaxConcurrent(3); got != 3 {
		t.Fatalf("resolveMaxConcurrent(3) = %d, want 3", got)
	}
}

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
		if err := RunCommandParallel(context.Background(), "pwd", nil, []ParallelDir{
			{Name: "repo-a", Path: dirA},
			{Name: "repo-b", Path: dirB},
		}, 0); err != nil {
			t.Errorf("RunCommandParallel returned unexpected error: %v", err)
		}
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
	// Deterministic order: repo-a before repo-b.
	if strings.Index(out, "REPO-A") > strings.Index(out, "REPO-B") {
		t.Fatalf("expected REPO-A before REPO-B in output, got:\n%s", out)
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
		if err := RunCommandParallel(context.Background(), "ls", []string{"hello.txt"}, []ParallelDir{
			{Name: "repo-a", Path: dirA},
		}, 0); err != nil {
			t.Errorf("RunCommandParallel returned unexpected error: %v", err)
		}
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
		if err := RunCommandParallel(context.Background(), "pwd", nil, []ParallelDir{
			{Name: "repo-ok", Path: dirOK},
			{Name: "repo-missing", Path: filepath.Join(tmp, "does-not-exist")},
		}, 0); err != nil {
			t.Errorf("RunCommandParallel returned unexpected error: %v", err)
		}
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
	if err := RunCommandParallel(context.Background(), "pwd", nil, dirs, 1); err != nil {
		t.Fatalf("RunCommandParallel returned unexpected error: %v", err)
	}
}

func TestRunCommandParallel_EmptyDirs(t *testing.T) {
	if err := RunCommandParallel(context.Background(), "pwd", nil, nil, 0); err != nil {
		t.Fatalf("RunCommandParallel returned unexpected error: %v", err)
	}
}

func TestRunCommandParallel_RespectsConcurrencyLimit(t *testing.T) {
	tmp := t.TempDir()
	script := filepath.Join(tmp, "work.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 0.2\n"), 0755); err != nil {
		t.Fatal(err)
	}

	dirs := make([]ParallelDir, 4)
	for i := range dirs {
		name := "d" + string(rune('a'+i))
		path := filepath.Join(tmp, name)
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
		dirs[i] = ParallelDir{Name: name, Path: path}
	}

	start := time.Now()
	if err := RunCommandParallel(context.Background(), script, nil, dirs, 1); err != nil {
		t.Fatalf("RunCommandParallel returned unexpected error: %v", err)
	}
	elapsed := time.Since(start)
	// 4 * 0.2s = 0.8s at limit 1; allow some slack for scheduling jitter.
	if elapsed < 700*time.Millisecond {
		t.Fatalf("expected ~sequential runtime with limit 1, got %v", elapsed)
	}
}

func TestRunCommandParallel_ContextCancel(t *testing.T) {
	tmp := t.TempDir()
	script := filepath.Join(tmp, "sleep.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\nsleep 5\n"), 0755); err != nil {
		t.Fatal(err)
	}

	var dirs []ParallelDir
	for _, name := range []string{"a", "b"} {
		path := filepath.Join(tmp, name)
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
		dirs = append(dirs, ParallelDir{Name: name, Path: path})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := RunCommandParallel(ctx, script, nil, dirs, 2)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
	if !strings.Contains(err.Error(), "context") {
		t.Fatalf("expected context error, got: %v", err)
	}
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
