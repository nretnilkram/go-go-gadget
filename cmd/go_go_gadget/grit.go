package go_go_gadget

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

func newGritCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grit",
		Short: "Run git commands on multiple repositories",
		Long: `Utility that allows you to run a git command on multiple git repository directories at once.

e.g. go-go-gadget grit pull

Will update all the of the repositories in the configuration.  Useful for updating all repositories in the morning.

Configuration (config.yml):
  max_concurrent  Maximum number of repositories to process in parallel.
                  Default: min(GOMAXPROCS*2, 16). Takes precedence over GRIT_MAX_CONCURRENT.

Environment Variables:
  GRIT_MAX_CONCURRENT  Maximum number of repositories to process in parallel.
                       Default: min(GOMAXPROCS*2, 16). Overridden by max_concurrent in config.yml when set.`,
		Args: cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := grit.TestGritDir(); err != nil {
				return err
			}
			return grit.AppendHistory(cmd.CommandPath() + " " + strings.Join(args, " "))
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
			defer stop()

			if err := grit.RunGitCommandParallel(ctx, args); err != nil {
				return err
			}
			writeln(cmd.OutOrStdout(), "Finished Run @ "+utilities.ShowDateTime("dash", true))
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}

	cmd.AddCommand(newGritAddAllReposCmd())
	cmd.AddCommand(newGritAddRepoCmd())
	cmd.AddCommand(newGritConfigCmd())
	cmd.AddCommand(newGritDestroyCmd())
	cmd.AddCommand(newGritHistoryCmd())
	cmd.AddCommand(newGritInitCmd())
	cmd.AddCommand(newGritRemoveRepoCmd())
	cmd.AddCommand(newGritResetCmd())
	return cmd
}

func newGritAddRepoCmd() *cobra.Command {
	return &cobra.Command{
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
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritAddAllReposCmd() *cobra.Command {
	return &cobra.Command{
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
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritConfigCmd() *cobra.Command {
	return &cobra.Command{
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

			out := cmd.OutOrStdout()
			writeln(out, "--------")
			writeln(out, string(yamlData))
			writeln(out, "--------")
			writeln(out, "")
			writeln(out, "Repositories Count: "+strconv.Itoa(len(config.Repositories)))
			writeln(out, "Grit Directory: "+cwd+"/"+grit.GritDir)
			writeln(out, "Config File: "+cwd+"/"+grit.ConfigFile)
			writeln(out, "History File: "+cwd+"/"+grit.HistoryFile)
			writeln(out, "Working Directory: "+cwd)

			grit.PrintTagLine(out, cmd.Root().Version)
			return nil
		},
	}
}

func newGritInitCmd() *cobra.Command {
	return &cobra.Command{
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
				writeln(cmd.OutOrStdout(), "Grit is already initialized.")
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

			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritHistoryCmd() *cobra.Command {
	return &cobra.Command{
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
			writeln(cmd.OutOrStdout(), string(history))
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritRemoveRepoCmd() *cobra.Command {
	return &cobra.Command{
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
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritResetCmd() *cobra.Command {
	return &cobra.Command{
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
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}

func newGritDestroyCmd() *cobra.Command {
	return &cobra.Command{
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
			grit.PrintTagLine(cmd.OutOrStdout(), cmd.Root().Version)
			return nil
		},
	}
}
