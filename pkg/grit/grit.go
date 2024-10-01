package grit

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var GritDir = ".grit"
var ConfigFile = GritDir + "/config.yml"
var HisotryFile = GritDir + "/history.log"

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func AppendHistory(command string) {
	// Open the file in append mode
	file, err := os.OpenFile(HisotryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err)
	defer file.Close()

	// Write the data to the file
	if _, err := file.WriteString(command + "\n"); err != nil {
		log.Fatal(err)
	}
}

func RunCommand(args []string) {
	fmt.Println("testing...")
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		fmt.Println(name)
		fmt.Println(path)
		fmt.Println("git " + strings.Join(args, " "))
	}
}
