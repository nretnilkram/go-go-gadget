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
	Short: "Run git commands on multiple repositories at once",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		if gritSyncronous {
			grit.RunGitCommandSyncronous(args)
		} else {
			grit.RunGitCommandParallel(args)
		}
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritAddRepoCmd = &cobra.Command{
	Use:     "add-repo",
	Aliases: []string{"add"},
	Short:   "Add repository to grit config",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.AddRepoToConfig(args[0], args[0])
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritAddAllReposCmd = &cobra.Command{
	Use:     "add-all-repos",
	Aliases: []string{"add-all"},
	Short:   "Add all git repositories in directory to grit config",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.AddAllRepos()
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show the grit conifig for the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		var config grit.Config = grit.LoadConfig()

		// Marshal the data into YAML format with indentation
		yamlData, err := yaml.Marshal(config)
		grit.Check(err)

		fmt.Println(string(yamlData))
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new Grit directory",
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
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show grit history",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		history, err := os.ReadFile(grit.HisotryFile)
		grit.Check(err)

		fmt.Println(string(history))
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritRemoveRepoCmd = &cobra.Command{
	Use:     "remove-repo",
	Aliases: []string{"remove"},
	Short:   "Remove repository from config",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.RemoveRepoFromConfig(args[0])
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset grit to default config",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		var config grit.Config = grit.DefaultConfig()
		grit.WriteConfig(config)

		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Remove grit from directory",
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		err := os.Remove(grit.GritDir)
		grit.Check(err)

		grit.PrintTagLine(cmd.Root().Version)
	},
}

func init() {
	gritCmd.AddCommand(gritAddAllReposCmd)

	gritCmd.AddCommand(gritAddRepoCmd)

	gritCmd.AddCommand(gritConfigCmd)

	gritCmd.AddCommand(gritDestroyCmd)

	gritCmd.AddCommand(gritHistoryCmd)

	gritCmd.AddCommand(gritInitCmd)

	gritCmd.AddCommand(gritRemoveRepoCmd)

	gritCmd.AddCommand(gritResetCmd)

	gritCmd.Flags().BoolVarP(&gritSyncronous, "syncronous", "s", false, "Run Grit Command Syncronously")
	rootCmd.AddCommand(gritCmd)
}
