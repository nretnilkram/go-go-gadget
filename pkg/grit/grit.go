package grit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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
	config := LoadConfig()
	failed := runGitCommandParallel(config, args)
	printFailedRepos(failed)
}

// RunGitCommandSynchronous runs the given git command sequentially across all configured repositories.
func RunGitCommandSynchronous(args []string) {
	config := LoadConfig()
	failed := runGitCommandSynchronous(config, args)
	printFailedRepos(failed)
}

func runGitCommandParallel(config Config, args []string) []string {
	maxConcurrent := 0
	if maxConcurrentStr := os.Getenv("GRIT_MAX_CONCURRENT"); maxConcurrentStr != "" {
		if parsed, err := strconv.Atoi(maxConcurrentStr); err == nil && parsed > 0 {
			maxConcurrent = parsed
		}
	}

	var wg sync.WaitGroup
	var failedMu sync.Mutex
	var failed []string

	var semaphore chan struct{}
	if maxConcurrent > 0 {
		semaphore = make(chan struct{}, maxConcurrent)
	}

	for _, repo := range config.Repositories {
		wg.Add(1)
		go func(repo Repository) {
			defer wg.Done()

			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}

			if runGitInRepo(config, repo, args, false) {
				failedMu.Lock()
				failed = append(failed, repo.Name)
				failedMu.Unlock()
			}
		}(repo)
	}
	wg.Wait()

	sort.Strings(failed)
	return failed
}

func runGitCommandSynchronous(config Config, args []string) []string {
	var failed []string

	for _, repo := range config.Repositories {
		if runGitInRepo(config, repo, args, true) {
			failed = append(failed, repo.Name)
		}
	}

	return failed
}

func runGitInRepo(config Config, repo Repository, args []string, synchronous bool) bool {
	commandDisplay := "git " + strings.Join(args, " ")
	repoDir := config.Root + "/" + repo.Path
	output, err := utilities.RunCommandWithError("git", args, repoDir)

	name := strings.ToUpper(repo.Name)
	if synchronous {
		fmt.Println(Header(name+" -- "+commandDisplay) + "\n" + output + Footer())
	} else {
		fmt.Println(Header(name+" -- ["+commandDisplay+"]") + "\n\n" + output + "\n" + Footer(name))
	}

	return err != nil
}

func printFailedRepos(failed []string) {
	if len(failed) == 0 {
		return
	}

	fmt.Println("\n" + Header("REPOSITORIES WITH ERRORS"))
	for _, name := range failed {
		fmt.Println("  " + name)
	}
	fmt.Print(Footer() + "\n")
}
