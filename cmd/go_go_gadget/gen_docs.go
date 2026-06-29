package go_go_gadget

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var genDocsCmd = &cobra.Command{
	Use:     "generate-documentation [md|rest|yaml]",
	Aliases: []string{"documentation", "docs"},
	Short:   "Generate Documentation",
	Long: `Generate Documentation for the Go Go Gadget CLI in MarkDown Rest or YAML format.

Aliases: documentation, docs`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	ValidArgs:             []string{"md", "rest", "yaml"},
	RunE: func(cmd *cobra.Command, args []string) error {
		doesDocsDirExist, err := utilities.FileDirExists("./docs")
		if err != nil {
			return err
		}

		if !doesDocsDirExist {
			if err := os.Mkdir("./docs", 0755); err != nil {
				return err
			}
		}

		switch args[0] {
		case "md":
			return doc.GenMarkdownTree(rootCmd, "./docs")
		case "rest":
			return doc.GenReSTTree(rootCmd, "./docs")
		case "yaml":
			return doc.GenYamlTree(rootCmd, "./docs")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genDocsCmd)
}
