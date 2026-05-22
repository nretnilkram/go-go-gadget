package grit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

// GritDir is the name of the grit metadata directory.
const GritDir = ".grit"

// ConfigFile is the path to the grit configuration file within GritDir.
const ConfigFile = GritDir + "/config.yml"

// HistoryFile is the path to the grit command history log within GritDir.
const HistoryFile = GritDir + "/history.log"

// AppendHistory appends a timestamped command entry to the grit history log.
func AppendHistory(command string) {
	// Open the file in append mode
	file, err := os.OpenFile(HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	utilities.Check(err)
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	// Write the data to the file
	if _, err := file.WriteString("[" + utilities.ShowDateTime("dash", true) + "] " + command + "\n"); err != nil {
		log.Fatal(err)
	}
}

// AddAllRepos scans the current directory for git repositories and adds any not already in config.
func AddAllRepos() {
	entries, err := os.ReadDir(".")
	utilities.Check(err)

	// Load config once to check for existing repos
	config := LoadConfig()
	existingRepos := make(map[string]bool)
	for _, repo := range config.Repositories {
		existingRepos[repo.Name] = true
		existingRepos[repo.Path] = true
	}

	// Loop over all directories and add to config if Git repository
	for _, entry := range entries {
		dotGitExists, _ := utilities.FileDirExists(filepath.Join(entry.Name() + "/.git"))
		if entry.IsDir() && dotGitExists {
			gitDir := entry.Name()
			// Only add if not already in configuration
			if !existingRepos[gitDir] {
				AddRepoToConfig(gitDir, gitDir)
			}
		}
	}
}

// RunGitCommandParallel runs the given git command concurrently across all configured repositories.
func RunGitCommandParallel(args []string) {
	// Create a map to store the parsed YAML data
	var config = LoadConfig()

	// Check for max concurrent limit from environment variable
	maxConcurrent := 0
	if maxConcurrentStr := os.Getenv("GRIT_MAX_CONCURRENT"); maxConcurrentStr != "" {
		if parsed, err := strconv.Atoi(maxConcurrentStr); err == nil && parsed > 0 {
			maxConcurrent = parsed
		}
	}

	// Create WaitGroup for parallel runs in all repositories
	var wg sync.WaitGroup

	// Create semaphore channel if maxConcurrent is set
	var semaphore chan struct{}
	if maxConcurrent > 0 {
		semaphore = make(chan struct{}, maxConcurrent)
	}

	for _, repo := range config.Repositories {
		wg.Add(1)
		go func(repo Repository) {
			defer wg.Done()

			// Acquire semaphore if concurrency limit is set
			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}

			path := repo.Path
			name := repo.Name
			// Use RunCommand to safely pass args without shell interpretation, preventing command injection
			commandDisplay := "git " + strings.Join(args, " ")
			repoDir := config.Root + "/" + path
			output := utilities.RunCommand("git", args, repoDir)
			fmt.Println(Header(strings.ToUpper(name)+" -- ["+commandDisplay+"]") + "\n\n" + output + "\n" + Footer(strings.ToUpper(name)))
		}(repo)
	}
	wg.Wait()
}

// RunGitCommandSynchronous runs the given git command sequentially across all configured repositories.
func RunGitCommandSynchronous(args []string) {
	// Create a map to store the parsed YAML data
	var config = LoadConfig()

	// Run command in all repositories
	for _, repo := range config.Repositories {
		path := repo.Path
		name := repo.Name
		// Use RunCommand to safely pass args without shell interpretation, preventing command injection
		commandDisplay := "git " + strings.Join(args, " ")
		repoDir := config.Root + "/" + path
		output := utilities.RunCommand("git", args, repoDir)
		fmt.Println(Header(strings.ToUpper(name)+" -- "+commandDisplay) + "\n" + output + Footer())
	}

}
