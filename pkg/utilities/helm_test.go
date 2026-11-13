package utilities

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBumpSemver(t *testing.T) {
	cases := []struct {
		version  string
		bumpType string
		want     string
		wantErr  bool
	}{
		{"1.2.3", "patch", "1.2.4", false},
		{"1.2.3", "minor", "1.3.0", false},
		{"1.2.3", "major", "2.0.0", false},
		{"1.2.3-beta.1", "patch", "1.2.4", false},
		{"1.2.3+build", "minor", "1.3.0", false},
		{"0.0.9", "PATCH", "0.0.10", false},
		{"1.0", "patch", "", true},
		{"1.2.3", "bogus", "", true},
	}
	for _, c := range cases {
		got, err := BumpSemver(c.version, c.bumpType)
		if c.wantErr {
			if err == nil {
				t.Errorf("BumpSemver(%q, %q) expected error", c.version, c.bumpType)
			}
			continue
		}
		if err != nil {
			t.Errorf("BumpSemver(%q, %q) unexpected error: %v", c.version, c.bumpType, err)
			continue
		}
		if got != c.want {
			t.Errorf("BumpSemver(%q, %q) = %q, want %q", c.version, c.bumpType, got, c.want)
		}
	}
}

func TestFindChartYAMLFiles(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "charts", "app")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	chartPath := filepath.Join(nested, "Chart.yaml")
	if err := os.WriteFile(chartPath, []byte("version: 0.1.0\n"), 0600); err != nil {
		t.Fatalf("write chart: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("x"), 0600); err != nil {
		t.Fatalf("write readme: %v", err)
	}

	got, err := FindChartYAMLFiles(root)
	if err != nil {
		t.Fatalf("FindChartYAMLFiles: %v", err)
	}
	if len(got) != 1 || got[0] != chartPath {
		t.Errorf("FindChartYAMLFiles = %v, want [%s]", got, chartPath)
	}
}

func TestBumpChartsInDir(t *testing.T) {
	root := t.TempDir()
	chartA := filepath.Join(root, "a")
	chartB := filepath.Join(root, "b")
	if err := os.MkdirAll(chartA, 0755); err != nil {
		t.Fatalf("mkdir a: %v", err)
	}
	if err := os.MkdirAll(chartB, 0755); err != nil {
		t.Fatalf("mkdir b: %v", err)
	}

	if err := os.WriteFile(filepath.Join(chartA, "Chart.yaml"), []byte("apiVersion: v2\nname: a\nversion: 1.2.3\n"), 0600); err != nil {
		t.Fatalf("write a: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chartB, "Chart.yaml"), []byte("apiVersion: v2\nname: b\nversion: \"2.0.0\"\n"), 0600); err != nil {
		t.Fatalf("write b: %v", err)
	}

	results, err := BumpChartsInDir(root, "minor")
	if err != nil {
		t.Fatalf("BumpChartsInDir: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	updated := 0
	for _, result := range results {
		if result.Skipped {
			t.Errorf("unexpected skip for %s: %s", result.File, result.SkipReason)
			continue
		}
		updated++
	}
	if updated != 2 {
		t.Fatalf("expected 2 updates, got %d", updated)
	}

	versionA, err := ReadChartVersion(filepath.Join(chartA, "Chart.yaml"))
	if err != nil {
		t.Fatalf("ReadChartVersion a: %v", err)
	}
	if versionA != "1.3.0" {
		t.Errorf("chart a version = %q, want 1.3.0", versionA)
	}

	contentB, err := os.ReadFile(filepath.Join(chartB, "Chart.yaml"))
	if err != nil {
		t.Fatalf("read b: %v", err)
	}
	if got := string(contentB); got != "apiVersion: v2\nname: b\nversion: \"2.1.0\"\n" {
		t.Errorf("chart b content = %q, want quoted 2.1.0 preserved", got)
	}
}

func TestBumpChartsInDir_SkipsMissingVersion(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "Chart.yaml"), []byte("apiVersion: v2\nname: missing\n"), 0600); err != nil {
		t.Fatalf("write chart: %v", err)
	}

	results, err := BumpChartsInDir(root, "patch")
	if err != nil {
		t.Fatalf("BumpChartsInDir: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Fatalf("expected skipped result, got %+v", results)
	}
}

func TestBumpChartsInDir_NoCharts(t *testing.T) {
	_, err := BumpChartsInDir(t.TempDir(), "patch")
	if err == nil {
		t.Fatal("expected error when no Chart.yaml files exist")
	}
}

func TestBumpChartsInDir_InvalidBumpType(t *testing.T) {
	_, err := BumpChartsInDir(t.TempDir(), "nope")
	if err == nil {
		t.Fatal("expected error for invalid bump type")
	}
}
