package go_go_gadget

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "2026.5.22"

var rootCmd = &cobra.Command{
	Use:     "go-go-gadget",
	Version: version,
	Short:   "go-go-gadget - CLI with useful tools",
	Long: `Go Go Gadget is a set of helpful CLI tools

You can use go-go-gadget to run useful tools from the terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\nGo Go Gadget...Help!\n\n* * * * * * * * * * * * *\n ")
		return cmd.Help()
	},
}

// Execute runs the root cobra command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'\n", err)
		os.Exit(1)
	}
}
