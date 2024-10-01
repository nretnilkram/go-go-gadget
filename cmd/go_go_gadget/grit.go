package go_go_gadget

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var gritCmd = &cobra.Command{
	Use:   "grit",
	Short: "Run git command on multiple repositories",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.RunCommand(args)
	},
}

var gritConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Run git command on multiple repositories",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		var config grit.Config = grit.LoadConfig()

		// Marshal the data into YAML format with indentation
		yamlData, err := yaml.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(yamlData))
	},
}

var gritInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init new Grit directory",
	Run: func(cmd *cobra.Command, args []string) {
		configFileExists, _ := grit.FileDirExists(grit.ConfigFile)
		if configFileExists {
			fmt.Println("Grit is already initialized.")
			return
		}

		var config grit.Config = grit.LoadConfig()

		f, err := os.Create(grit.HisotryFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		grit.WriteConfig(config)
	},
}

func init() {
	gritCmd.AddCommand(gritConfigCmd)

	gritCmd.AddCommand(gritInitCmd)

	rootCmd.AddCommand(gritCmd)
}
