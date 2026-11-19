package grit

import (
	"context"
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
// New repositories are appended in memory and written once.
func AddAllRepos() error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read directory: %w", err)
	}

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	existingRepos := make(map[string]bool, len(config.Repositories)*2)
	for _, repo := range config.Repositories {
		existingRepos[repo.Name] = true
		existingRepos[repo.Path] = true
	}

	added := false
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		gitDir := entry.Name()
		dotGitExists, err := utilities.FileDirExists(filepath.Join(gitDir, ".git"))
		if err != nil {
			return fmt.Errorf("check %s: %w", gitDir, err)
		}
		if !dotGitExists || existingRepos[gitDir] {
			continue
		}

		fmt.Println("Adding " + gitDir)
		config.Repositories = append(config.Repositories, Repository{
			Name: gitDir,
			Path: gitDir,
		})
		existingRepos[gitDir] = true
		added = true
	}

	if !added {
		return nil
	}
	return WriteConfig(config)
}

// RunGitCommandParallel runs the given git command concurrently across all configured repositories.
func RunGitCommandParallel(ctx context.Context, args []string) error {
	return RunCommandParallel(ctx, "git", args)
}

// RunCommandParallel runs the given command concurrently across all configured repositories.
// When max_concurrent / GRIT_MAX_CONCURRENT are unset or <= 0, utilities.DefaultMaxConcurrent is used.
func RunCommandParallel(ctx context.Context, commandName string, args []string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	dirs := make([]utilities.ParallelDir, len(config.Repositories))
	for i, repo := range config.Repositories {
		dirs[i] = utilities.ParallelDir{
			Name: repo.Name,
			Path: config.Root + "/" + repo.Path,
		}
	}

	return utilities.RunCommandParallel(ctx, commandName, args, dirs, maxConcurrent(config))
}

func maxConcurrent(config Config) int {
	limit := 0
	if maxConcurrentStr := os.Getenv("GRIT_MAX_CONCURRENT"); maxConcurrentStr != "" {
		if parsed, err := strconv.Atoi(maxConcurrentStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	// Config value takes precedence over the environment variable.
	if config.MaxConcurrent > 0 {
		limit = config.MaxConcurrent
	}
	return limit
}
