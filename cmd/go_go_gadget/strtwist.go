package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/go-go-gadget/pkg/strtwist"
	"github.com/spf13/cobra"
)

var onlyDigits bool
var inspectCmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"insp"},
	Short:   "Inspects a string",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		i := args[0]
		res, kind := strtwist.Inspect(i, onlyDigits)

		pluralS := "s"
		if res == 1 {
			pluralS = ""
		}
		fmt.Printf("'%s' has %d %s%s.\n", i, res, kind, pluralS)
	},
}

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetesify a string",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := strtwist.K8s(args[0])
		fmt.Println(res)
	},
}

var reverseCmd = &cobra.Command{
	Use:     "reverse",
	Aliases: []string{"rev"},
	Short:   "Reverses a string",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := strtwist.Reverse(args[0])
		fmt.Println(res)
	},
}

var symSubCmd = &cobra.Command{
	Use:     "symsub",
	Aliases: []string{"sym", "sub"},
	Short:   "Substitute symbols into a string",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := strtwist.SymbolSubstitution(args[0])
		fmt.Println(res)
	},
}

func init() {
	inspectCmd.Flags().BoolVarP(&onlyDigits, "digits", "d", false, "Count only digits")
	rootCmd.AddCommand(inspectCmd)

	rootCmd.AddCommand(k8sCmd)

	rootCmd.AddCommand(symSubCmd)

	rootCmd.AddCommand(reverseCmd)
}
