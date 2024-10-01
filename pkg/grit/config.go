package grit

import (
	"os"

	"gopkg.in/yaml.v2"
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

func DefaultConfig() Config {
	config := Config{
		Root:        "asdf",
		Ignore_Root: true,
	}

	return config
}

func LoadConfig() Config {
	// Read the file content
	data, err := os.ReadFile(ConfigFile)
	Check(err)

	// Create a map to store the parsed YAML data
	var config Config

	// Unmarshal the YAML string into the map
	err = yaml.Unmarshal([]byte(data), &config)
	Check(err)

	return config
}

func WriteConfig(config Config) {
	// Marshal the data into YAML format with indentation
	yamlData, err := yaml.Marshal(config)
	Check(err)

	data := []byte("---\n" + string(yamlData))

	writeErr := os.WriteFile(ConfigFile, data, 0644)
	Check(writeErr)
}

func AddRepoToConfig(name string, path string) {
	config := LoadConfig()
	repo := Repository{
		Name: name,
		Path: path,
	}
	config.Repositories = append(config.Repositories, repo)
	WriteConfig(config)
}
