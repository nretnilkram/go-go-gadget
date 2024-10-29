package grit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var GritDir = ".grit"
var ConfigFile = GritDir + "/config.yml"
var HistoryFile = GritDir + "/history.log"

func AppendHistory(command string) {
	// Open the file in append mode
	file, err := os.OpenFile(HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utilities.Check(err)
	defer file.Close()

	// Create timestamp
	now := time.Now().Format("[2006-01-02 15:04:05]")

	// Write the data to the file
	if _, err := file.WriteString(now + " " + command + "\n"); err != nil {
		log.Fatal(err)
	}
}

func AddAllRepos() {
	entries, err := os.ReadDir(".")
	utilities.Check(err)

	// Loop over all directories and add to config
	for _, entry := range entries {
		if entry.IsDir() && utilities.IsGitRepo(filepath.Join(entry.Name())) {
			gitDir := entry.Name()
			AddRepoToConfig(gitDir, gitDir)
		}
	}
}

func RunGitCommandParallel(args []string) {
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	// Create WaitGroup for parallel runs in all repositories
	var wg sync.WaitGroup
	for _, repo := range config.Repositories {
		wg.Add(1)
		go func(repo Repository) {
			defer wg.Done()
			path := repo.Path
			name := repo.Name
			command := "git " + strings.Join(args, " ")
			repoDir := config.Root + "/" + path
			fmt.Println(GritHeader(strings.ToUpper(name)+" -- "+command) + "\n" + utilities.RunShellCommand(command, repoDir) + GritFooter())
		}(repo)
	}
	wg.Wait()
}

func RunGitCommandSynchronous(args []string) {
	// Create a map to store the parsed YAML data
	var config Config = LoadConfig()

	// Run command in all repositories
	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		command := "git " + strings.Join(args, " ")
		repoDir := config.Root + "/" + path
		fmt.Println(GritHeader(strings.ToUpper(name)+" -- "+command) + "\n" + utilities.RunShellCommand(command, repoDir) + GritFooter())
	}

}
