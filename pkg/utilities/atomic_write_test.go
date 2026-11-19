package utilities

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileAtomic_CreatesAndReplaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")

	if err := WriteFileAtomic(path, []byte("first\n"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic create: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "first\n" {
		t.Fatalf("content = %q, want %q", got, "first\n")
	}

	if err := WriteFileAtomic(path, []byte("second\n"), 0644); err != nil {
		t.Fatalf("WriteFileAtomic replace: %v", err)
	}
	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "second\n" {
		t.Fatalf("content = %q, want %q", got, "second\n")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected only final file, found %d entries", len(entries))
	}
}
