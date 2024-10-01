package go_go_gadget

import (
	"fmt"
	"log"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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

		var config map[string]interface{} = grit.LoadConfig()

		// Marshal the data into YAML format with indentation
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(yamlData))
	},
}

func init() {
	gritCmd.AddCommand(gritConfigCmd)

	rootCmd.AddCommand(gritCmd)
}
