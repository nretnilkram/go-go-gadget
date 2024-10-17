package go_go_gadget

import (
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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
	Run: func(cmd *cobra.Command, args []string) {
		doesDocsDirExist, err := utilities.FileDirExists("./docs")
		utilities.Check(err)

		if !(doesDocsDirExist) {
			dirErr := os.Mkdir("./docs", 0755)
			utilities.Check(dirErr)
		}

		switch args[0] {
		case "md":
			genErr := doc.GenMarkdownTree(rootCmd, "./docs")
			utilities.Check(genErr)
		case "rest":
			genErr := doc.GenReSTTree(rootCmd, "./docs")
			utilities.Check(genErr)
		case "yaml":
			genErr := doc.GenYamlTree(rootCmd, "./docs")
			utilities.Check(genErr)
		}
	},
}

func init() {
	rootCmd.AddCommand(genDocsCmd)
}
