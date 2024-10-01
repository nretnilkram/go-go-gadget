package grit

import (
	"fmt"
	"os"
)

// Check the error status
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Get the Current Working Directory
func GetWorkingDir() string {
	dir, err := os.Getwd()
	Check(err)
	return dir
}

var WorkingDir = GetWorkingDir()

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

func GritHeader(header string) string {
	return "----------------------------------------\n# " + header
}

func GritFooter() string {
	return "----------------------------------------\n"
}

func PrintHeader(header string) {
	fmt.Println(GritHeader((header)))
}

func PrintFooter() {
	fmt.Println(GritFooter())
}
