package grit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

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

	// Write the data to the file
	if _, err := file.WriteString("[" + utilities.ShowDateTime("dash", true) + "] " + command + "\n"); err != nil {
		log.Fatal(err)
	}
}

func AddAllRepos() {
	entries, err := os.ReadDir(".")
	utilities.Check(err)

	// Loop over all directories and add to config if Git repository
	for _, entry := range entries {
		dotGitExists, _ := utilities.FileDirExists(filepath.Join(entry.Name() + "/.git"))
		if entry.IsDir() && dotGitExists {
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
