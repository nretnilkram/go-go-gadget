package grit

import (
	"fmt"
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
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
		Root:        utilities.GetWorkingDir(),
		Ignore_Root: true,
	}

	return config
}

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

func WriteConfig(config Config) {
	// Marshal the data into YAML format with indentation
	yamlData, err := yaml.Marshal(config)
	utilities.Check(err)

	data := []byte("---\n" + string(yamlData))

	writeErr := os.WriteFile(ConfigFile, data, 0644)
	utilities.Check(writeErr)
}

func AddRepoToConfig(name string, path string) {
	config := LoadConfig()
	repo := Repository{
		Name: name,
		Path: path,
	}
	fmt.Println("Adding " + name)
	config.Repositories = append(config.Repositories, repo)
	WriteConfig(config)
}

func RemoveRepoFromConfig(name string) {
	config := LoadConfig()

	for i, repo := range config.Repositories {
		if name == repo.Name {
			config.Repositories = append(config.Repositories[:i], config.Repositories[i+1:]...)
			break
		}
	}
	fmt.Println("Removing " + name)
	WriteConfig(config)
}
