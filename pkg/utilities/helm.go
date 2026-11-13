package utilities

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	bumpableSemverRe = regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9]+)(?:[-+].*)?$`)
	chartVersionRe   = regexp.MustCompile(`^version:[[:space:]]*["']?([^"'[:space:]]+)["']?[[:space:]]*$`)
)

// ChartBumpResult describes one Chart.yaml version bump attempt.
type ChartBumpResult struct {
	File       string
	OldVersion string
	NewVersion string
	Skipped    bool
	SkipReason string
}

// BumpSemver increments the major, minor, or patch segment of a version string.
// Pre-release and build suffixes are dropped, matching common chart bump behavior.
func BumpSemver(version, bumpType string) (string, error) {
	matches := bumpableSemverRe.FindStringSubmatch(version)
	if matches == nil {
		return "", fmt.Errorf("unsupported version format: %s", version)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return "", err
	}
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", err
	}
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return "", err
	}

	switch strings.ToLower(bumpType) {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return "", fmt.Errorf("invalid bump type: %s", bumpType)
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// FindChartYAMLFiles returns sorted paths to Chart.yaml files under root.
func FindChartYAMLFiles(root string) ([]string, error) {
	var charts []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == "Chart.yaml" {
			charts = append(charts, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return charts, nil
}

// ReadChartVersion extracts the version field from a Chart.yaml file.
func ReadChartVersion(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return extractChartVersion(string(content))
}

func extractChartVersion(content string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		matches := chartVersionRe.FindStringSubmatch(line)
		if matches != nil {
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("no version field found")
}

func replaceChartVersion(content, currentVersion, newVersion string) (string, error) {
	lines := strings.Split(content, "\n")
	updated := false
	for i, line := range lines {
		matches := chartVersionRe.FindStringSubmatch(line)
		if matches == nil || matches[1] != currentVersion {
			continue
		}
		prefix := line[:strings.Index(line, matches[1])]
		suffix := line[strings.Index(line, matches[1])+len(matches[1]):]
		lines[i] = prefix + newVersion + suffix
		updated = true
		break
	}
	if !updated {
		return "", fmt.Errorf("version line not updated")
	}
	return strings.Join(lines, "\n"), nil
}

// BumpChartFile bumps the version in a single Chart.yaml and writes it back.
func BumpChartFile(path, bumpType string) (ChartBumpResult, error) {
	result := ChartBumpResult{File: path}

	content, err := os.ReadFile(path)
	if err != nil {
		return result, err
	}

	currentVersion, err := extractChartVersion(string(content))
	if err != nil {
		result.Skipped = true
		result.SkipReason = err.Error()
		return result, nil
	}
	result.OldVersion = currentVersion

	newVersion, err := BumpSemver(currentVersion, bumpType)
	if err != nil {
		return result, err
	}
	result.NewVersion = newVersion

	updated, err := replaceChartVersion(string(content), currentVersion, newVersion)
	if err != nil {
		return result, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return result, err
	}
	if err := os.WriteFile(path, []byte(updated), info.Mode().Perm()); err != nil {
		return result, err
	}

	return result, nil
}

// BumpChartsInDir finds Chart.yaml files under dir and bumps each version.
func BumpChartsInDir(dir, bumpType string) ([]ChartBumpResult, error) {
	bumpType = strings.ToLower(bumpType)
	switch bumpType {
	case "major", "minor", "patch":
	default:
		return nil, fmt.Errorf("invalid bump type: %s", bumpType)
	}

	charts, err := FindChartYAMLFiles(dir)
	if err != nil {
		return nil, err
	}
	if len(charts) == 0 {
		return nil, fmt.Errorf("no Chart.yaml files found under %s", dir)
	}

	results := make([]ChartBumpResult, 0, len(charts))
	for _, chart := range charts {
		result, err := BumpChartFile(chart, bumpType)
		if err != nil {
			return results, fmt.Errorf("%s: %w", chart, err)
		}
		results = append(results, result)
	}
	return results, nil
}
