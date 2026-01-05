package grit

import (
	"fmt"
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var WorkingDir = utilities.GetWorkingDir()

// Check if directory contains a grit dir and config
func TestGritDir() {
	gritDirExists, _ := utilities.FileDirExists(GritDir)
	configFileExists, _ := utilities.FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		fmt.Println("This is not a grit directory.")
		os.Exit(1)
	}
}

// Add header of dashes with optional header string
func GritHeader(headerString ...string) string {
	if len(headerString) > 0 {
		return "----------------------------------------\n>> " + headerString[0]
	}
	return "----------------------------------------\n"
}

// Add footer of dashes with optional footer string
func GritFooter(footerString ...string) string {
	if len(footerString) > 0 {
		return "<< " + footerString[0] + "\n" + "----------------------------------------\n"
	}
	return "----------------------------------------\n"
}

// Print Header with optional header string
func PrintHeader(headerString ...string) {
	fmt.Println(GritHeader(headerString...))
}

// Print Footer with optional footer string
func PrintFooter(footerString ...string) {
	fmt.Println(GritFooter(footerString...))
}

// Print version tag line
func PrintTagLine(version string) {
	fmt.Println("\n--------\ngo-go-gadget grit v" + version + "\n")
}
