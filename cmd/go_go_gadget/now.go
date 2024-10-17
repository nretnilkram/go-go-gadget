package go_go_gadget

import (
	"fmt"
	"time"

	"github.com/nretnilkram/go-go-gadget/pkg/now"
	"github.com/spf13/cobra"
)

var useColons bool
var useDots bool
var useGoRawTime bool
var useSlashes bool
var useUnix bool
var includeTime bool

var nowCmd = &cobra.Command{
	Use:     "now",
	Aliases: []string{"today"},
	Short:   "Print todays date",
	Long: `Print out todays date with a number of formatting options including separators and time.

Aliases: now, today`,
	Run: func(cmd *cobra.Command, args []string) {
		if useUnix {
			fmt.Println(time.Now().Unix())
		} else if useGoRawTime {
			fmt.Println(now.ShowDateTime("raw", includeTime))
		} else if useColons {
			fmt.Println(now.ShowDateTime("colon", includeTime))
		} else if useDots {
			fmt.Println(now.ShowDateTime("dot", includeTime))
		} else if useSlashes {
			fmt.Println(now.ShowDateTime("slash", includeTime))
		} else {
			fmt.Println(now.ShowDateTime("dash", includeTime))
		}
	},
}

func init() {
	nowCmd.Flags().BoolVarP(&useColons, "colon", "c", false, "Use colons")
	nowCmd.Flags().BoolVarP(&useDots, "dots", "d", false, "Use dots")
	nowCmd.Flags().BoolVarP(&useGoRawTime, "raw", "r", false, "Raw go time")
	nowCmd.Flags().BoolVarP(&useSlashes, "slashes", "s", false, "Use Slases")
	nowCmd.Flags().BoolVarP(&useUnix, "unix", "u", false, "Include time")
	nowCmd.Flags().BoolVarP(&includeTime, "time", "t", false, "Include time")
	rootCmd.AddCommand(nowCmd)
}
