package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/pswd"
	"github.com/spf13/cobra"
)

var passwordLength int
var passwordSymbols bool

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"p", "pw"},
	Short:   "Generate a password",
	Long: `Generate a password with customizable length and complexity.

Aliases: password, p, pw`,
	Run: func(cmd *cobra.Command, args []string) {
		weight := pswd.PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 0}

		if passwordSymbols {
			weight = pswd.PasswordWeight{Lower: 4, Upper: 3, Digit: 3, Symbol: 2}
		}

		res := pswd.Password(passwordLength, weight)
		fmt.Println(res)
	},
}

func init() {
	passwordCmd.Flags().BoolVarP(&passwordSymbols, "symbols", "s", true, "Include Symbols")
	passwordCmd.Flags().IntVarP(&passwordLength, "length", "l", 32, "Length of Password")
	rootCmd.AddCommand(passwordCmd)
}
