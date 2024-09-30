package go_go_gadget

import (
	"fmt"
	"os"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/spf13/cobra"
)

var gritCmd = &cobra.Command{
	Use:   "grit",
	Short: "Run git command on multiple repositories",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		grit.RunCommand(args)
	},
}

var gritConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Run git command on multiple repositories",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		// Read the file content
		data, err := os.ReadFile(".grit/config.yml")
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		// Convert the byte slice to a string
		content := string(data)

		fmt.Println(content)
	},
}

func init() {
	gritCmd.AddCommand(gritConfigCmd)

	rootCmd.AddCommand(gritCmd)
}
