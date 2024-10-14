package utilities

import (
	"bytes"
	"os"
	"os/exec"
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
		return "ERROR: " + string(err.Error()) + "\n" + stderr.String()
	}

	return out.String()
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
