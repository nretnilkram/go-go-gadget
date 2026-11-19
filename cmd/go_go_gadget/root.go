package go_go_gadget

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var version = "2026.11.19"

// newRootCmd builds a fresh command tree. Tests should call this rather than
// sharing package-level command state.
func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "go-go-gadget",
		Version: version,
		Short:   "go-go-gadget - CLI with useful tools",
		Long: `Go Go Gadget is a set of helpful CLI tools

You can use go-go-gadget to run useful tools from the terminal.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			writeln(cmd.OutOrStdout(), "\nGo Go Gadget...Help!\n\n* * * * * * * * * * * * *\n ")
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(newInspectCmd())
	rootCmd.AddCommand(newK8sCmd())
	rootCmd.AddCommand(newReverseCmd())
	rootCmd.AddCommand(newSymSubCmd())
	rootCmd.AddCommand(newNowCmd())
	rootCmd.AddCommand(newWordsCmd())
	rootCmd.AddCommand(newPasswordCmd())
	rootCmd.AddCommand(newM8sCmd())
	rootCmd.AddCommand(newUtilCmd())
	rootCmd.AddCommand(newGritCmd())
	rootCmd.AddCommand(newGroupCmd())
	rootCmd.AddCommand(newGenDocsCmd())

	return rootCmd
}

// Execute runs the root cobra command and exits on error.
func Execute() {
	ExecuteWith(os.Stdout, os.Stderr)
}

// ExecuteWith runs the CLI using the given writers. Intended for tests and
// programmatic embedding; the process still exits on error to match Execute.
func ExecuteWith(stdout, stderr io.Writer) {
	root := newRootCmd()
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		writef(stderr, "Whoops. There was an error while executing your CLI '%s'\n", err)
		os.Exit(1)
	}
}
