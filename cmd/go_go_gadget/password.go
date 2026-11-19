package go_go_gadget

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/password"
)

var passwordLength int
var passwordSymbols bool

var passwordCmd = &cobra.Command{
	Use:     "password",
	Aliases: []string{"p", "pw"},
	Short:   "Generate a password",
	Long: `Generate a password with customizable length and complexity.

Use --symbols=false to generate a password without symbols.

Aliases: password, p, pw`,
	RunE: func(cmd *cobra.Command, args []string) error {
		weight := password.Weight{Lower: 4, Upper: 3, Digit: 3, Symbol: 0}
		if passwordSymbols {
			weight.Symbol = 2
		}

		res, err := password.Generate(passwordLength, weight)
		if err != nil {
			return err
		}
		fmt.Println(res)
		return nil
	},
}

func init() {
	passwordCmd.Flags().BoolVarP(&passwordSymbols, "symbols", "s", true, "Include Symbols")
	passwordCmd.Flags().IntVarP(&passwordLength, "length", "l", 32, "Length of Password")
	rootCmd.AddCommand(passwordCmd)
}
