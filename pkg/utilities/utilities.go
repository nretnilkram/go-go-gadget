package utilities

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// SemverRegex is a regular expression pattern for validating semantic version strings.
var SemverRegex = "^(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"

// Check the error status
func Check(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", e)
		os.Exit(1)
	}
}

// RunCommand runs a command with the given name and arguments safely, without shell interpretation.
// This prevents command injection by passing arguments directly to exec.Command.
// Returns the combined stdout and stderr output as a string.
func RunCommand(commandName string, args []string, path string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(commandName, args...)
	cmd.Dir = path
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return string(err.Error()) + "\n" + stderr.String() + "\n" + out.String()
	}

	return out.String() + stderr.String()
}

// RunCommandInteractive runs a command with the given name and arguments safely, with shell interpretation and interactive input.
func RunCommandInteractive(commandName string, args []string, path string) {

	cmd := exec.Command(commandName, args...)
	cmd.Dir = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(string(err.Error()))
	}
}

// FileDirExists returns whether the given file or directory exists
func FileDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsGitRepo returns true only when the given path is inside a git repository.
// It checks the exit code rather than scanning output, which avoids false
// positives from error messages that themselves contain ".git".
func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	return cmd.Run() == nil
}

// GetWorkingDir returns the current working directory, panicking on error.
func GetWorkingDir() string {
	dir, err := os.Getwd()
	Check(err)
	return dir
}

// GrepFileForTFResources scans a Terraform file and returns a slice of --target flags
// for all resource and module blocks found.
func GrepFileForTFResources(filename string) []string {
	pattern := "^(resource |module)"

	resources := []string{}

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		os.Exit(1)
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			line = strings.ReplaceAll(line, "resource \"", "--target=")
			line = strings.ReplaceAll(line, "module \"", "--target=module.")
			line = strings.ReplaceAll(line, "\" \"", ".")
			line = strings.ReplaceAll(line, "\" {", " \\")
			resources = append(resources, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return (resources)
}

// ListTFResources prints the Terraform resource targets found in each of the given files.
func ListTFResources(files []string) {
	for _, file := range files {
		fmt.Println(strings.Join(GrepFileForTFResources(file), "\n"))
	}
}

// RegexTest reports whether input matches the given regular expression pattern.
func RegexTest(input string, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		os.Exit(1)
	}

	return re.MatchString(input)
}

// WaitForConfirmationPrompt displays a [y/n] prompt and returns true if the user confirms.
func WaitForConfirmationPrompt(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]: ", s)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	if response == "y" || response == "Y" {
		return true
	}
	return false
}

// ShowDateTime returns the current date/time as a formatted string using the given separator style.
// Pass showTime as true to include the time component.
func ShowDateTime(format string, showTime bool) (dateTime string) {
	currentTime := time.Now()
	var year = currentTime.Year()
	var month = int(currentTime.Month())
	var day = currentTime.Day()
	var hour = currentTime.Hour()
	var minute = currentTime.Minute()
	var second = currentTime.Second()
	var separator string

	switch format {
	case "colon":
		separator = ":"
	case "dash":
		separator = "-"
	case "dot":
		separator = "."
	case "slash":
		separator = "/"
	default:
		return currentTime.String()
	}

	dateTime = fmt.Sprintf("%d%s%d%s%d", year, separator, month, separator, day)
	if showTime {
		dateTime = dateTime + " " + fmt.Sprintf("%d:%02d:%02d", hour, minute, second)
	}

	return dateTime
}
