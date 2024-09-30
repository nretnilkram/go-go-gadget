package go_go_gadget

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "2.0.0"

var rootCmd = &cobra.Command{
	Use:     "go_go_gadget",
	Version: version,
	Short:   "go_go_gadget - a simple CLI with useful tools",
	Long: `Go Go Gadget is a super fancy CLI (kidding)

One can use go_go_gadget to run useful tools from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nGo Go Gadget...Help!\n\n* * * * * * * * * * * * *\n ")
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
