package go_go_gadget

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var gritSynchronous bool

var gritCmd = &cobra.Command{
	Use:   "grit",
	Short: "Run git commands on multiple repositories",
	Long: `Utility that allows you to run a git command on multiple git repository directories at once.

e.g. go-go-gadget grit pull

Will update all the of the repositories in the configuration.  Useful for updating all repositories in the morning.

Environment Variables:
  GRIT_MAX_CONCURRENT  Maximum number of repositories to process in parallel (default: unlimited).`,
	Args: cobra.MinimumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := grit.TestGritDir(); err != nil {
			return err
		}
		return grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if gritSynchronous {
			grit.RunGitCommandSynchronous(args)
		} else {
			grit.RunGitCommandParallel(args)
		}
		fmt.Println("Finished Run @ " + utilities.ShowDateTime("dash", true))
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritAddRepoCmd = &cobra.Command{
	Use:     "add-repo",
	Aliases: []string{"add"},
	Short:   "Add repository",
	Long: `Add a new repository to your grit configuration.

Aliases: add-repo, add`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := grit.AddRepoToConfig(args[0], args[0]); err != nil {
			return err
		}
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritAddAllReposCmd = &cobra.Command{
	Use:     "add-all-repos",
	Aliases: []string{"add-all"},
	Short:   "Add all repositories",
	Long: `Add all git repositories in directory to grit config.

Aliases: add-all-repos, add-all`,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := grit.AddAllRepos(); err != nil {
			return err
		}
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritConfigCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"conf"},
	Short:   "Show config",
	Long: `Print the current grit configuration.

Aliases: config, conf`,
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := grit.LoadConfig()
		if err != nil {
			return err
		}

		yamlData, err := yaml.Marshal(config)
		if err != nil {
			return err
		}

		cwd, err := utilities.GetWorkingDir()
		if err != nil {
			return err
		}

		fmt.Println("--------")
		fmt.Println(string(yamlData))
		fmt.Println("--------")
		fmt.Println("")
		fmt.Println("Repositories Count: " + strconv.Itoa(len(config.Repositories)))
		fmt.Println("Grit Directory: " + cwd + "/" + grit.GritDir)
		fmt.Println("Config File: " + cwd + "/" + grit.ConfigFile)
		fmt.Println("History File: " + cwd + "/" + grit.HistoryFile)
		fmt.Println("Working Directory: " + cwd)

		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritInitCmd = &cobra.Command{
	Use:     "initialize",
	Aliases: []string{"init"},
	Short:   "Initialize Grit",
	Long: `Initialize current directory with a .grit directory and new config file.

Aliases: initialize, init`,
	DisableFlagsInUseLine: true,
	PersistentPreRunE:     func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		configFileExists, _ := utilities.FileDirExists(grit.GritDir)
		if configFileExists {
			fmt.Println("Grit is already initialized.")
			return nil
		}

		if err := os.Mkdir(grit.GritDir, 0755); err != nil {
			return err
		}

		config, err := grit.DefaultConfig()
		if err != nil {
			return err
		}
		if err := grit.WriteConfig(config); err != nil {
			return err
		}

		f, err := os.Create(grit.HistoryFile)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Println(err)
			}
		}()

		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritHistoryCmd = &cobra.Command{
	Use:                   "history",
	Short:                 "Show grit history",
	Long:                  "Print the history of the current grit directory.",
	DisableFlagsInUseLine: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return grit.TestGritDir()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		history, err := os.ReadFile(grit.HistoryFile)
		if err != nil {
			return err
		}
		fmt.Println(string(history))
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritRemoveRepoCmd = &cobra.Command{
	Use:     "remove-repo",
	Aliases: []string{"remove", "rm"},
	Short:   "Remove repository",
	Long: `Remove repositories from your grit configuration.

The argument is a glob pattern matched against repository names.
Use % as the wildcard (e.g. gg-phoenix-%) because shells expand unquoted *
before grit runs. Quoted * also works (e.g. 'gg-phoenix-*').

Aliases: remove-repo, remove, rm`,
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return grit.TestGritDir()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := grit.RemoveRepoFromConfig(args[0]); err != nil {
			return err
		}
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritResetCmd = &cobra.Command{
	Use:                   "reset",
	Short:                 "Reset grit",
	Long:                  "Reset grit configuration to the default configuration.",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if utilities.WaitForConfirmationPrompt("Do you want to continue?") {
			config, err := grit.DefaultConfig()
			if err != nil {
				return err
			}
			if err := grit.WriteConfig(config); err != nil {
				return err
			}
		}
		grit.PrintTagLine(cmd.Root().Version)
		return nil
	},
}

var gritDestroyCmd = &cobra.Command{
	Use:                   "destroy",
	Short:                 "Clean grit",
	Long:                  "Cleanup the current grit setup by removing the .grit directory and contents.",
	DisableFlagsInUseLine: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return grit.TestGritDir()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.RemoveAll(grit.GritDir); err != nil {
			return err
		}
		grit.PrintTagLine(cmd.Root().Version)
		return nil
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
