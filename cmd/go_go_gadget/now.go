package go_go_gadget

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

func newNowCmd() *cobra.Command {
	var (
		useColons    bool
		useDots      bool
		useGoRawTime bool
		useSlashes   bool
		useUnix      bool
		includeTime  bool
	)

	cmd := &cobra.Command{
		Use:     "now",
		Aliases: []string{"today"},
		Short:   "Print todays date",
		Long: `Print out todays date with a number of formatting options including separators and time.

Aliases: now, today`,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			switch {
			case useUnix:
				writeln(out, time.Now().Unix())
			case useGoRawTime:
				writeln(out, utilities.ShowDateTime("raw", includeTime))
			case useColons:
				writeln(out, utilities.ShowDateTime("colon", includeTime))
			case useDots:
				writeln(out, utilities.ShowDateTime("dot", includeTime))
			case useSlashes:
				writeln(out, utilities.ShowDateTime("slash", includeTime))
			default:
				writeln(out, utilities.ShowDateTime("dash", includeTime))
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&useColons, "colon", "c", false, "Use colons")
	cmd.Flags().BoolVarP(&useDots, "dots", "d", false, "Use dots")
	cmd.Flags().BoolVarP(&useGoRawTime, "raw", "r", false, "Raw go time")
	cmd.Flags().BoolVarP(&useSlashes, "slashes", "s", false, "Use Slashes")
	cmd.Flags().BoolVarP(&useUnix, "unix", "u", false, "Print Unix timestamp")
	cmd.Flags().BoolVarP(&includeTime, "time", "t", false, "Include time")
	return cmd
}
