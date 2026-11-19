package go_go_gadget

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/grit"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

var groupCmd = &cobra.Command{
	Use:   "group [command] [args...]",
	Short: "Run a command on multiple directories (uses grit config to determine directories)",
	Long: `Run any command in parallel across all directories listed in the grit configuration.

Uses the same .grit/config.yml and concurrency settings as grit
(max_concurrent / GRIT_MAX_CONCURRENT; default min(GOMAXPROCS*2, 16)).

e.g. go-go-gadget group ls -la
e.g. go-go-gadget group npm install
e.g. go-go-gadget group make test`,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return grit.TestGritDir()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer stop()

		if err := grit.RunCommandParallel(ctx, args[0], args[1:]); err != nil {
			return err
		}
		fmt.Println("Finished Run @ " + utilities.ShowDateTime("dash", true))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
}
