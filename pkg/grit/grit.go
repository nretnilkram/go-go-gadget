package grit

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

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
	RunCommandParallel("git", args)
}

// RunCommandParallel runs the given command concurrently across all configured repositories.
func RunCommandParallel(commandName string, args []string) {
	config, err := LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	dirs := make([]utilities.ParallelDir, len(config.Repositories))
	for i, repo := range config.Repositories {
		dirs[i] = utilities.ParallelDir{
			Name: repo.Name,
			Path: config.Root + "/" + repo.Path,
		}
	}

	utilities.RunCommandParallel(commandName, args, dirs, maxConcurrent(config))
}

func maxConcurrent(config Config) int {
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
	return maxConcurrent
}
