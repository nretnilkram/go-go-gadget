package utilities

import (
	"context"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// maxConcurrentCap caps the automatic default worker count.
const maxConcurrentCap = 16

// ParallelDir is a named directory target for parallel command execution.
type ParallelDir struct {
	Name string
	Path string
}

type dirFailure struct {
	Name   string
	Output string
}

type dirResult struct {
	display string
	failure *dirFailure
}

// DefaultMaxConcurrent returns a bounded worker count based on GOMAXPROCS.
// Values are clamped to [1, maxConcurrentCap].
func DefaultMaxConcurrent() int {
	n := runtime.GOMAXPROCS(0) * 2
	if n > maxConcurrentCap {
		return maxConcurrentCap
	}
	if n < 1 {
		return 1
	}
	return n
}

func resolveMaxConcurrent(maxConcurrent int) int {
	if maxConcurrent <= 0 {
		return DefaultMaxConcurrent()
	}
	return maxConcurrent
}

// RunCommandParallel runs the given command concurrently in each directory.
// maxConcurrent limits parallelism; values <= 0 use DefaultMaxConcurrent().
// Returns ctx.Err() when the context is cancelled before all work finishes.
func RunCommandParallel(ctx context.Context, commandName string, args []string, dirs []ParallelDir, maxConcurrent int) error {
	failed, err := runCommandParallel(ctx, commandName, args, dirs, maxConcurrent)
	printFailedDirs(failed)
	return err
}

func buildParallelDirOutput(ctx context.Context, dir ParallelDir, commandName string, args []string) dirResult {
	commandDisplay := commandName
	if len(args) > 0 {
		commandDisplay += " " + strings.Join(args, " ")
	}
	output, err := RunCommandWithErrorContext(ctx, commandName, args, dir.Path)
	name := strings.ToUpper(dir.Name)
	display := parallelHeader(name+" -- ["+commandDisplay+"]") + "\n\n" + output + "\n" + parallelFooter(name)
	result := dirResult{display: display}
	if err != nil {
		result.failure = &dirFailure{Name: dir.Name, Output: output}
	}
	return result
}

func runCommandParallel(ctx context.Context, commandName string, args []string, dirs []ParallelDir, maxConcurrent int) ([]dirFailure, error) {
	if len(dirs) == 0 {
		return nil, nil
	}

	limit := resolveMaxConcurrent(maxConcurrent)
	if limit > len(dirs) {
		limit = len(dirs)
	}

	results := make([]dirResult, len(dirs))
	jobs := make(chan int)

	var wg sync.WaitGroup
	for range limit {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				if ctx.Err() != nil {
					return
				}
				results[i] = buildParallelDirOutput(ctx, dirs[i], commandName, args)
			}
		}()
	}

sendLoop:
	for i := range dirs {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()

	var failed []dirFailure
	for _, result := range results {
		if result.display == "" {
			continue
		}
		fmt.Print(result.display)
		if result.failure != nil {
			failed = append(failed, *result.failure)
		}
	}

	sort.Slice(failed, func(i, j int) bool {
		return failed[i].Name < failed[j].Name
	})

	if err := ctx.Err(); err != nil {
		return failed, err
	}
	return failed, nil
}

func printFailedDirs(failed []dirFailure) {
	if len(failed) == 0 {
		return
	}

	fmt.Println(parallelHeader("DIRECTORIES WITH ERRORS"))
	for _, f := range failed {
		fmt.Println("  " + f.Name + ":")
		for _, line := range strings.Split(strings.TrimRight(f.Output, "\n"), "\n") {
			fmt.Println("    " + line)
		}
	}
	fmt.Print(parallelFooter() + "\n")
}

func parallelHeader(headerString ...string) string {
	if len(headerString) > 0 {
		return "----------------------------------------\n>> " + headerString[0]
	}
	return "----------------------------------------\n"
}

func parallelFooter(footerString ...string) string {
	if len(footerString) > 0 {
		return "<< " + footerString[0] + "\n" + "----------------------------------------\n"
	}
	return "----------------------------------------\n"
}
