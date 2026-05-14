package grit

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

type Repository struct {
	Name string
	Path string
}

type Config struct {
	Root         string
	Repositories []Repository
	Ignore_Root  bool
}

// Generate and return Default Config
func DefaultConfig() Config {
	config := Config{
		Root:        utilities.GetWorkingDir(),
		Ignore_Root: true,
	}

	return config
}

// Load and return config from filesystem
func LoadConfig() Config {
	// Read the file content
	data, err := os.ReadFile(ConfigFile)
	utilities.Check(err)

	// Create a map to store the parsed YAML data
	var config Config

	// Unmarshal the YAML string into the map
	err = yaml.Unmarshal([]byte(data), &config)
	utilities.Check(err)

	return config
}

// Save passed config to filesystem
func WriteConfig(config Config) {
	// Marshal the data into YAML format with indentation
	yamlData, err := yaml.Marshal(config)
	utilities.Check(err)

	data := []byte("---\n" + string(yamlData))

	writeErr := os.WriteFile(ConfigFile, data, 0644)
	utilities.Check(writeErr)
}

// Add passed repository to configuration
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

// Remove passed repository from configuration
func RemoveRepoFromConfig(name string) {
	config := LoadConfig()

	found := false
	for i, repo := range config.Repositories {
		if name == repo.Name {
			config.Repositories = append(config.Repositories[:i], config.Repositories[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Repository " + name + " not found in configuration.")
		return
	}

	fmt.Println("Removing " + name)
	WriteConfig(config)
}
