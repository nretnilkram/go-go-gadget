package grit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Repository represents a single git repository entry in the grit config.
type Repository struct {
	Name string
	Path string
}

// Config holds the root path and list of repositories managed by grit.
type Config struct {
	Root          string
	Repositories  []Repository
	IgnoreRoot    bool `yaml:"ignore_root"`
	MaxConcurrent int  `yaml:"max_concurrent"`
}

// DefaultConfig returns a Config populated with default values.
func DefaultConfig() (Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("get working directory: %w", err)
	}
	return Config{
		Root:       dir,
		IgnoreRoot: true,
	}, nil
}

// LoadConfig reads and unmarshals the grit config file from the filesystem.
func LoadConfig() (Config, error) {
	data, err := os.ReadFile(ConfigFile)
	if err != nil {
		return Config{}, fmt.Errorf("read %s: %w", ConfigFile, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", ConfigFile, err)
	}

	return config, nil
}

// WriteConfig marshals the given config and writes it to the grit config file.
func WriteConfig(config Config) error {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	data := []byte("---\n" + string(yamlData))
	if err := os.WriteFile(ConfigFile, data, 0644); err != nil {
		return fmt.Errorf("write %s: %w", ConfigFile, err)
	}
	return nil
}

// AddRepoToConfig adds a repository with the given name and path to the grit config.
func AddRepoToConfig(name string, path string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	for _, repo := range config.Repositories {
		if name == repo.Name || path == repo.Path {
			fmt.Println("Repository " + name + " already exists in configuration.")
			return nil
		}
	}

	repo := Repository{
		Name: name,
		Path: path,
	}
	fmt.Println("Adding " + name)
	config.Repositories = append(config.Repositories, repo)
	return WriteConfig(config)
}

// normalizeRemovePattern converts shell-safe % wildcards to * for filepath.Match.
// Use % instead of * on the command line so zsh/bash do not expand the pattern first.
func normalizeRemovePattern(pattern string) string {
	return strings.ReplaceAll(pattern, "%", "*")
}

// RemoveRepoFromConfig removes repositories whose names match the given pattern from the grit config.
// The pattern uses glob syntax. Prefer % over * on the command line (e.g. gg-phoenix-%)
// because shells expand unquoted * before grit receives the argument.
func RemoveRepoFromConfig(pattern string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	glob := normalizeRemovePattern(pattern)

	var removed []string
	var kept []Repository
	for _, repo := range config.Repositories {
		matched, err := filepath.Match(glob, repo.Name)
		if err != nil {
			fmt.Println("Invalid pattern " + pattern + ": " + err.Error())
			return nil
		}
		if matched {
			removed = append(removed, repo.Name)
		} else {
			kept = append(kept, repo)
		}
	}

	if len(removed) == 0 {
		fmt.Println("Repository " + pattern + " not found in configuration.")
		return nil
	}

	for _, name := range removed {
		fmt.Println("Removing " + name)
	}
	config.Repositories = kept
	return WriteConfig(config)
}
