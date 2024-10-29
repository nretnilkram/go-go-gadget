package go_go_gadget

import (
	"fmt"
	"os"
	"strings"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var gritSynchronous bool

var gritCmd = &cobra.Command{
	Use:   "grit",
	Short: "Run git commands on multiple repositories",
	Long: `Utility that allows you to run a git command on multiple git repository directories at once.

e.g. go-go-gadget grit pull

Will update all the of the repositories in the configuration.  Useful for updating all repositories in the morning.`,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		if gritSynchronous {
			grit.RunGitCommandSynchronous(args)
		} else {
			grit.RunGitCommandParallel(args)
		}
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritAddRepoCmd = &cobra.Command{
	Use:     "add-repo",
	Aliases: []string{"add"},
	Short:   "Add repository",
	Long: `Add a new repository to your grit configuration.

Aliases: app-repo, add`,
	DisableFlagsInUseLine: true,
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
	Short:   "Add all repositories",
	Long: `Add all git repositories in directory to grit config.

Aliases: add-all-repos, add-all`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.AddAllRepos()
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Show config",
	Long: `Print the current grig configuration.

Aliases: config, conf`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		var config grit.Config = grit.LoadConfig()

		// Marshal the data into YAML format with indentation
		yamlData, err := yaml.Marshal(config)
		utilities.Check(err)

		fmt.Println(string(yamlData))
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritInitCmd = &cobra.Command{
	Use:     "initialize",
	Aliases: []string{"init"},
	Short:   "Initialize Grit",
	Long: `Initialize current directory with a .grit directory and new config file.

Aliases: initialize, init`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		configFileExists, _ := utilities.FileDirExists(grit.GritDir)
		if configFileExists {
			fmt.Println("Grit is already initialized.")
			return
		}

		// Create .grit Dir
		dirErr := os.Mkdir(grit.GritDir, 0755)
		utilities.Check(dirErr)

		// Create Default Config File
		var config grit.Config = grit.DefaultConfig()
		grit.WriteConfig(config)

		// Create Hisotry File
		f, historyErr := os.Create(grit.HistoryFile)
		utilities.Check(historyErr)
		defer f.Close()
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritHistoryCmd = &cobra.Command{
	Use:                   "history",
	Short:                 "Show grit history",
	Long:                  "Print the history of the current grit directory.",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		history, err := os.ReadFile(grit.HistoryFile)
		utilities.Check(err)

		fmt.Println(string(history))
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritRemoveRepoCmd = &cobra.Command{
	Use:     "remove-repo",
	Aliases: []string{"remove", "rm"},
	Short:   "Remove repository",
	Long: `Remove a new repository to your grit configuration.

Aliases: remove-repo, remove, rm`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		grit.RemoveRepoFromConfig(args[0])
		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritResetCmd = &cobra.Command{
	Use:                   "reset",
	Short:                 "Reset grit",
	Long:                  "Reset grit configuration to the default configuration.",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()
		grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))

		if utilities.WaitForConfirmationPrompt("Do you want to continue?") {
			var config grit.Config = grit.DefaultConfig()
			grit.WriteConfig(config)
		}

		grit.PrintTagLine(cmd.Root().Version)
	},
}

var gritDestroyCmd = &cobra.Command{
	Use:                   "destroy",
	Short:                 "Clean grit",
	Long:                  "Cleanup the current grit setup by removing the .grit directory and contents.",
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		grit.TestGritDir()

		err := os.Remove(grit.GritDir)
		utilities.Check(err)

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

	gritCmd.Flags().BoolVarP(&gritSynchronous, "synchronous", "s", false, "Run Grit Command Synchronously")
	rootCmd.AddCommand(gritCmd)
}
