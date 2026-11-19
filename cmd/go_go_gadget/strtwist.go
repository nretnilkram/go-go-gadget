package go_go_gadget

import (
	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/strtwist"
)

func newInspectCmd() *cobra.Command {
	var onlyDigits bool

	cmd := &cobra.Command{
		Use:     "inspect",
		Aliases: []string{"insp"},
		Short:   "Inspects a string",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			i := args[0]
			res, kind := strtwist.Inspect(i, onlyDigits)

			pluralS := "s"
			if res == 1 {
				pluralS = ""
			}
			writef(cmd.OutOrStdout(), "'%s' has %d %s%s.\n", i, res, kind, pluralS)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&onlyDigits, "digits", "d", false, "Count only digits")
	return cmd
}

func newK8sCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "k8s",
		Short: "Kubernetesify a string",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			writeln(cmd.OutOrStdout(), strtwist.K8s(args[0]))
			return nil
		},
	}
}

func newReverseCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "reverse",
		Aliases: []string{"rev"},
		Short:   "Reverses a string",
		Long: `Takes a string a returns the characters in reverse.

Aliases: reverse, rev`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			writeln(cmd.OutOrStdout(), strtwist.Reverse(args[0]))
			return nil
		},
	}
}

func newSymSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "symsub",
		Aliases: []string{"sym", "sub"},
		Short:   "Substitute symbols",
		Long: `Takes a string and returns the same string with a subset of characters substituted symbols.

'a', 'A' -> '@'
'e', 'E' -> '3'
'i', 'I' -> '!'
'o', 'O' -> '0'
's', 'S' -> '$'

Aliases: symsub, sym, sub`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			writeln(cmd.OutOrStdout(), strtwist.SymbolSubstitution(args[0]))
			return nil
		},
	}
}
