package grit

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var GritDir = ".grit"
var ConfigFile = GritDir + "/config.yml"

func TestGritDir() {
	gritDirExists, _ := FileDirExists(GritDir)
	configFileExists, _ := FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		fmt.Println("This is not a grit directory.")
		os.Exit(1)
	}
}

// FileDirExists returns whether the given file or directory exists
func FileDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Repository struct {
	Name string
	Path string
}

type Config struct {
	Root         string
	Repositories []Repository
	Ingore_Root  bool
}

func LoadConfig() Config {
	// Read the file content
	data, err := os.ReadFile(".grit/config.yml")
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}

	// Create a map to store the parsed YAML data
	var config Config

	// Unmarshal the YAML string into the map
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func RunCommand(args []string) {
	fmt.Println("testing...")
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	for _, repo := range config.Repositories {
		path := repo.Path
		fmt.Println(path)
	}
}
