package grit

import (
	"fmt"
	"os"
)

// exists returns whether the given file or directory exists
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

func RunCommand(args []string) {
	fmt.Println("testing...")
}
