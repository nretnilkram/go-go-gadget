package utilities

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

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

// RunCommandParallel runs the given command concurrently in each directory.
// maxConcurrent limits parallelism; values <= 0 mean unlimited.
func RunCommandParallel(commandName string, args []string, dirs []ParallelDir, maxConcurrent int) {
	failed := runCommandParallel(commandName, args, dirs, maxConcurrent)
	printFailedDirs(failed)
}

func buildParallelDirOutput(dir ParallelDir, commandName string, args []string) dirResult {
	commandDisplay := commandName
	if len(args) > 0 {
		commandDisplay += " " + strings.Join(args, " ")
	}
	output, err := RunCommandWithError(commandName, args, dir.Path)
	name := strings.ToUpper(dir.Name)
	display := parallelHeader(name+" -- ["+commandDisplay+"]") + "\n\n" + output + "\n" + parallelFooter(name)
	result := dirResult{display: display}
	if err != nil {
		result.failure = &dirFailure{Name: dir.Name, Output: output}
	}
	return result
}

func runCommandParallel(commandName string, args []string, dirs []ParallelDir, maxConcurrent int) []dirFailure {
	var wg sync.WaitGroup
	var printMu sync.Mutex
	var failedMu sync.Mutex
	var failed []dirFailure

	var semaphore chan struct{}
	if maxConcurrent > 0 {
		semaphore = make(chan struct{}, maxConcurrent)
	}

	for _, dir := range dirs {
		wg.Add(1)
		go func(dir ParallelDir) {
			defer wg.Done()

			if semaphore != nil {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
			}

			result := buildParallelDirOutput(dir, commandName, args)

			// Print atomically so blocks from concurrent dirs never interleave.
			printMu.Lock()
			fmt.Print(result.display)
			printMu.Unlock()

			if result.failure != nil {
				failedMu.Lock()
				failed = append(failed, *result.failure)
				failedMu.Unlock()
			}
		}(dir)
	}

	wg.Wait()

	sort.Slice(failed, func(i, j int) bool {
		return failed[i].Name < failed[j].Name
	})
	return failed
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
