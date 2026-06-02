package grit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

// Repository represents a single git repository entry in the grit config.
type Repository struct {
	Name string
	Path string
}

// Config holds the root path and list of repositories managed by grit.
type Config struct {
	Root         string
	Repositories []Repository
	IgnoreRoot   bool `yaml:"ignore_root"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() Config {
	config := Config{
		Root:       utilities.GetWorkingDir(),
		IgnoreRoot: true,
	}

	return config
}

// LoadConfig reads and unmarshals the grit config file from the filesystem.
func LoadConfig() Config {
	// Read the file content
	data, err := os.ReadFile(ConfigFile)
	utilities.Check(err)

	// Create a map to store the parsed YAML data
	var config Config

	// Unmarshal the YAML string into the map
	err = yaml.Unmarshal(data, &config)
	utilities.Check(err)

	return config
}

// WriteConfig marshals the given config and writes it to the grit config file.
func WriteConfig(config Config) {
	// Marshal the data into YAML format with indentation
	yamlData, err := yaml.Marshal(config)
	utilities.Check(err)

	data := []byte("---\n" + string(yamlData))

	writeErr := os.WriteFile(ConfigFile, data, 0644)
	utilities.Check(writeErr)
}

// AddRepoToConfig adds a repository with the given name and path to the grit config.
func AddRepoToConfig(name string, path string) {
	config := LoadConfig()

	// Check if repository already exists
	for _, repo := range config.Repositories {
		if name == repo.Name || path == repo.Path {
			fmt.Println("Repository " + name + " already exists in configuration.")
			return
		}
	}

	repo := Repository{
		Name: name,
		Path: path,
	}
	fmt.Println("Adding " + name)
	config.Repositories = append(config.Repositories, repo)
	WriteConfig(config)
}

// normalizeRemovePattern converts shell-safe % wildcards to * for filepath.Match.
// Use % instead of * on the command line so zsh/bash do not expand the pattern first.
func normalizeRemovePattern(pattern string) string {
	return strings.ReplaceAll(pattern, "%", "*")
}

// RemoveRepoFromConfig removes repositories whose names match the given pattern from the grit config.
// The pattern uses glob syntax. Prefer % over * on the command line (e.g. gg-phoenix-%)
// because shells expand unquoted * before grit receives the argument.
func RemoveRepoFromConfig(pattern string) {
	config := LoadConfig()
	glob := normalizeRemovePattern(pattern)

	var removed []string
	var kept []Repository
	for _, repo := range config.Repositories {
		matched, err := filepath.Match(glob, repo.Name)
		if err != nil {
			fmt.Println("Invalid pattern " + pattern + ": " + err.Error())
			return
		}
		if matched {
			removed = append(removed, repo.Name)
		} else {
			kept = append(kept, repo)
		}
	}

	if len(removed) == 0 {
		fmt.Println("Repository " + pattern + " not found in configuration.")
		return
	}

	for _, name := range removed {
		fmt.Println("Removing " + name)
	}
	config.Repositories = kept
	WriteConfig(config)
}
