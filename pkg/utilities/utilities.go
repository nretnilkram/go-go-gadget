package utilities

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

// Check the error status
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

var lastQuote = rune(0)
var f = func(c rune) bool {
	switch {
	case c == lastQuote:
		lastQuote = rune(0)
		return false
	case lastQuote != rune(0):
		return false
	case unicode.In(c, unicode.Quotation_Mark):
		lastQuote = c
		return false
	default:
		return unicode.IsSpace(c)
	}
}

// Run Shell Command and return result as string
func RunShellCommand(command string, path string) string {
	commandArray := strings.FieldsFunc(command, f)

	var out bytes.Buffer
	var stderr bytes.Buffer

	app := commandArray[0]
	args := commandArray[1:]

	cmd := exec.Command(app, args...)
	cmd.Dir = path
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return string(err.Error()) + "\n" + stderr.String() + "\n" + out.String()
	}

	return out.String() + stderr.String()
}

// Run Shell Command and return result as string
func RunShellCommandv2(app string, args []string, path string) string {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(app, args...)
	cmd.Dir = path
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return string(err.Error()) + "\n" + stderr.String() + "\n" + out.String()
	}

	return out.String() + stderr.String()
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

// Test if current directory is a git repository
func IsGitRepo(path string) bool {
	doesGitDirExist, err := FileDirExists(path + "/.git")
	Check(err)
	return doesGitDirExist
}

// Get the Current Working Directory
func GetWorkingDir() string {
	dir, err := os.Getwd()
	Check(err)
	return dir
}

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
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			line = strings.Replace(line, "resource ", "--target=", -1)
			line = strings.Replace(line, "module \"", "--target=module.", -1)
			line = strings.Replace(line, "\" \"", ".", -1)
			line = strings.Replace(line, "\" {", " \\", -1)
			resources = append(resources, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	return (resources)
}

func ListTFResources(files []string) {
	for _, file := range files {
		fmt.Println(strings.Join(GrepFileForTFResources(file), "\n"))
	}
}

func RegexTest(input string, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		os.Exit(1)
	}

	return re.MatchString(input)
}
