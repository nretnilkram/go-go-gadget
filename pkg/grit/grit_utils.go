package grit

import (
	"errors"
	"fmt"
	"io"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

// GetWorkingDir returns the current working directory at the time of the call.
func GetWorkingDir() (string, error) {
	return utilities.GetWorkingDir()
}

// TestGritDir checks if the current directory contains a grit dir and config.
// Returns an error if either is missing.
func TestGritDir() error {
	gritDirExists, _ := utilities.FileDirExists(GritDir)
	configFileExists, _ := utilities.FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		workingDir, _ := GetWorkingDir()
		return errors.New(workingDir + " is not a grit directory")
	}
	return nil
}

// Header returns a dashed header line with an optional header string appended.
func Header(headerString ...string) string {
	if len(headerString) > 0 {
		return "----------------------------------------\n>> " + headerString[0]
	}
	return "----------------------------------------\n"
}

// Footer returns a dashed footer line with an optional footer string prepended.
func Footer(footerString ...string) string {
	if len(footerString) > 0 {
		return "<< " + footerString[0] + "\n" + "----------------------------------------\n"
	}
	return "----------------------------------------\n"
}

// PrintHeader writes a grit header line with an optional header string.
func PrintHeader(w io.Writer, headerString ...string) {
	_, _ = fmt.Fprintln(w, Header(headerString...))
}

// PrintFooter writes a grit footer line with an optional footer string.
func PrintFooter(w io.Writer, footerString ...string) {
	_, _ = fmt.Fprintln(w, Footer(footerString...))
}

// PrintTagLine writes the grit version tag line.
func PrintTagLine(w io.Writer, version string) {
	_, _ = fmt.Fprintln(w, "\n--------\ngo-go-gadget grit v"+version+"\n")
}
