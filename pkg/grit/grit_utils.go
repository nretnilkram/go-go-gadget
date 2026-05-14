package grit

import (
	"fmt"
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

// WorkingDir holds the working directory at process startup.
var WorkingDir = utilities.GetWorkingDir()

// TestGritDir checks if the current directory contains a grit dir and config, exiting if not.
func TestGritDir() {
	gritDirExists, _ := utilities.FileDirExists(GritDir)
	configFileExists, _ := utilities.FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		fmt.Println("This is not a grit directory.")
		os.Exit(1)
	}
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

// PrintHeader prints a grit header line with an optional header string.
func PrintHeader(headerString ...string) {
	fmt.Println(Header(headerString...))
}

// PrintFooter prints a grit footer line with an optional footer string.
func PrintFooter(footerString ...string) {
	fmt.Println(Footer(footerString...))
}

// PrintTagLine prints the grit version tag line.
func PrintTagLine(version string) {
	fmt.Println("\n--------\ngo-go-gadget grit v" + version + "\n")
}
