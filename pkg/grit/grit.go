package grit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var GritDir = ".grit"
var ConfigFile = GritDir + "/config.yml"
var HisotryFile = GritDir + "/history.log"

func TestGritDir() {
	gritDirExists, _ := FileDirExists(GritDir)
	configFileExists, _ := FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		fmt.Println("This is not a grit directory.")
		os.Exit(1)
	}
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

func RunShellCommand(command string, path string) {
	commandArray := strings.Split(command, " ")

	cmd := exec.Command(commandArray[0], strings.Join(commandArray[1:], " "))
	cmd.Dir = path

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println(string(output))
		return
	}

	fmt.Println(string(output))
}

func RunGitCommand(args []string) {
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		command := "git " + strings.Join(args, " ")
		repoDir := config.Root + "/" + path
		PrintHeader(name + " - " + command)
		RunShellCommand(command, repoDir)
		PrintFooter()
	}
}

func RunGitPullCommand(args []string) {
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		command := "git pull"
		repoDir := config.Root + "/" + path
		PrintHeader(name + " - " + command)
		RunShellCommand(command, repoDir)
		PrintFooter()
	}
}
