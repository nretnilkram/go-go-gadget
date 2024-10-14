package grit

import (
	"fmt"
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var WorkingDir = utilities.GetWorkingDir()

func TestGritDir() {
	gritDirExists, _ := utilities.FileDirExists(GritDir)
	configFileExists, _ := utilities.FileDirExists(ConfigFile)

	if !gritDirExists || !configFileExists {
		fmt.Println("This is not a grit directory.")
		os.Exit(1)
	}
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

func PrintTagLine(version string) {
	fmt.Println("\n--------\ngo-go-gadget grit v" + version + "\n")
}
