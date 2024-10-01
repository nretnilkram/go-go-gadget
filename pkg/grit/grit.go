package grit

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

	// Create timestamp
	now := time.Now().Format("[2006-01-02 15:04:05]")

	// Write the data to the file
	if _, err := file.WriteString(now + " " + command + "\n"); err != nil {
		log.Fatal(err)
	}
}

func isGitRepo(path string) bool {
	doesGitDirExist, err := FileDirExists(path + "/.git")
	Check(err)
	return doesGitDirExist
}

func AddAllRepos() {
	entries, err := os.ReadDir(".")
	Check(err)

	for _, entry := range entries {
		if entry.IsDir() && isGitRepo(filepath.Join(entry.Name())) {
			gitDir := entry.Name()
			AddRepoToConfig(gitDir, gitDir)
		}
	}
}

func RunShellCommand(command string, path string) string {
	commandArray := strings.Split(command, " ")

	cmd := exec.Command(commandArray[0], strings.Join(commandArray[1:], " "))
	cmd.Dir = path

	// Run the command and capture its output
	output, err := cmd.Output()
	if err != nil {
		errString := err.Error()
		return "ERROR: " + string(errString) + "\n"
	}

	return string(output)
}

func RunGitCommandParallel(args []string) {
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()
	// Create WaitGroup for parallel runs
	var wg sync.WaitGroup
	for _, repo := range config.Repositories {
		wg.Add(1)
		go func(repo Repository) {
			defer wg.Done()
			path := repo.Path
			name := repo.Name
			command := "git " + strings.Join(args, " ")
			repoDir := config.Root + "/" + path
			fmt.Println(GritHeader(strings.ToUpper(name)+" -- "+command) + "\n" + RunShellCommand(command, repoDir) + GritFooter())
		}(repo)
	}
	wg.Wait()
}

func RunGitCommandSyncronous(args []string) {

	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		command := "git " + strings.Join(args, " ")
		repoDir := config.Root + "/" + path
		fmt.Println(GritHeader(strings.ToUpper(name)+" -- "+command) + "\n" + RunShellCommand(command, repoDir) + GritFooter())
	}

}
