package grit

import (
	"fmt"
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
func AppendHistory(command string) error {
	file, err := os.OpenFile(HistoryFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open history file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "close history file: %v\n", closeErr)
		}
	}()

	if _, err := file.WriteString("[" + utilities.ShowDateTime("dash", true) + "] " + command + "\n"); err != nil {
		return fmt.Errorf("write history: %w", err)
	}
	return nil
}

// AddAllRepos scans the current directory for git repositories and adds any not already in config.
func AddAllRepos() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read directory: %w", err)
	}

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	existingRepos := make(map[string]bool)
	for _, repo := range config.Repositories {
		existingRepos[repo.Name] = true
		existingRepos[repo.Path] = true
	}

	for _, entry := range entries {
		dotGitExists, _ := utilities.FileDirExists(filepath.Join(entry.Name() + "/.git"))
		if entry.IsDir() && dotGitExists {
			gitDir := entry.Name()
			if !existingRepos[gitDir] {
				if err := AddRepoToConfig(gitDir, gitDir); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// RunGitCommandParallel runs the given git command concurrently across all configured repositories.
func RunGitCommandParallel(args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	failed := runGitCommandParallel(config, args)
	printFailedRepos(failed)
}

// RunGitCommandSynchronous runs the given git command sequentially across all configured repositories.
func RunGitCommandSynchronous(args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}
	failed := runGitCommandSynchronous(config, args)
	printFailedRepos(failed)
}

type repoFailure struct {
	Name   string
	Output string
}

// repoResult holds the display output and optional failure for one repo.
type repoResult struct {
	display string
	failure *repoFailure
}

func buildParallelRepoOutput(config Config, repo Repository, args []string) repoResult {
	commandDisplay := "git " + strings.Join(args, " ")
	repoDir := config.Root + "/" + repo.Path
	output, err := utilities.RunCommandWithError("git", args, repoDir)
	name := strings.ToUpper(repo.Name)
	display := Header(name+" -- ["+commandDisplay+"]") + "\n\n" + output + "\n" + Footer(name)
	result := repoResult{display: display}
	if err != nil {
		result.failure = &repoFailure{Name: repo.Name, Output: output}
	}
	return result
}

func runGitCommandParallel(config Config, args []string) []repoFailure {
	maxConcurrent := 0
	if maxConcurrentStr := os.Getenv("GRIT_MAX_CONCURRENT"); maxConcurrentStr != "" {
		if parsed, err := strconv.Atoi(maxConcurrentStr); err == nil && parsed > 0 {
			maxConcurrent = parsed
		}
	}
	// Config value takes precedence over the environment variable.
	if config.MaxConcurrent > 0 {
		maxConcurrent = config.MaxConcurrent
	}

	var wg sync.WaitGroup
	var printMu sync.Mutex
	var failedMu sync.Mutex
	var failed []repoFailure

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

			result := buildParallelRepoOutput(config, repo, args)

			// Print atomically so blocks from concurrent repos never interleave.
			printMu.Lock()
			fmt.Print(result.display)
			printMu.Unlock()

			if result.failure != nil {
				failedMu.Lock()
				failed = append(failed, *result.failure)
				failedMu.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	sort.Slice(failed, func(i, j int) bool {
		return failed[i].Name < failed[j].Name
	})
	return failed
}

func runGitCommandSynchronous(config Config, args []string) []repoFailure {
	var failed []repoFailure

	for _, repo := range config.Repositories {
		if failure := runGitInRepoSync(config, repo, args); failure != nil {
			failed = append(failed, *failure)
		}
	}

	return failed
}

func runGitInRepoSync(config Config, repo Repository, args []string) *repoFailure {
	commandDisplay := "git " + strings.Join(args, " ")
	repoDir := config.Root + "/" + repo.Path
	output, err := utilities.RunCommandWithError("git", args, repoDir)

	name := strings.ToUpper(repo.Name)
	fmt.Println(Header(name+" -- "+commandDisplay) + "\n" + output + Footer())

	if err != nil {
		return &repoFailure{Name: repo.Name, Output: output}
	}
	return nil
}

func printFailedRepos(failed []repoFailure) {
	if len(failed) == 0 {
		return
	}

	fmt.Println(Header("REPOSITORIES WITH ERRORS"))
	for _, f := range failed {
		fmt.Println("  " + f.Name + ":")
		for _, line := range strings.Split(strings.TrimRight(f.Output, "\n"), "\n") {
			fmt.Println("    " + line)
		}
	}
	fmt.Print(Footer() + "\n")
}
