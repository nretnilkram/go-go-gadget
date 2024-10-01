package go_go_gadget

import (
	"fmt"
	"os"
	"strings"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var gritSyncronous bool

var gritCmd = &cobra.Command{
	Use:   "grit",
	Short: "Run git command on multiple repositories",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		if gritSyncronous {
			grit.RunGitCommandSyncronous(args)
		} else {
			grit.RunGitCommandParallel(args)
		}
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
		grit.Check(err)

		fmt.Println(string(yamlData))
	},
}

var gritInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Init new Grit directory",
	Run: func(cmd *cobra.Command, args []string) {
		configFileExists, _ := grit.FileDirExists(grit.GritDir)
		if configFileExists {
			fmt.Println("Grit is already initialized.")
			return
		}

		// Create .grit Dir
		dirErr := os.Mkdir(grit.GritDir, 0755)
		grit.Check(dirErr)

		// Create Default Config File
		var config grit.Config = grit.DefaultConfig()
		grit.WriteConfig(config)

		// Create Hisotry File
		f, historyErr := os.Create(grit.HisotryFile)
		grit.Check(historyErr)
		defer f.Close()
	},
}

var gritResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove grit from directory",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
	},
}

func init() {
	gritCmd.AddCommand(gritConfigCmd)

	gritCmd.AddCommand(gritInitCmd)

	gritCmd.AddCommand(gritResetCmd)

	gritCmd.Flags().BoolVarP(&gritSyncronous, "syncronous", "s", false, "Run Grit Command Syncronously")
	rootCmd.AddCommand(gritCmd)
}
