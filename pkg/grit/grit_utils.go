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

// Add header of dashes with header string
func GritHeader(header string) string {
	return "----------------------------------------\n# " + header
}

// Add footer of dashes
func GritFooter() string {
	return "----------------------------------------\n"
}

// Print Header
func PrintHeader(header string) {
	fmt.Println(GritHeader((header)))
}

// Print Footer
func PrintFooter() {
	fmt.Println(GritFooter())
}

// Print version tag line
func PrintTagLine(version string) {
	fmt.Println("\n--------\ngo-go-gadget grit v" + version + "\n")
}
