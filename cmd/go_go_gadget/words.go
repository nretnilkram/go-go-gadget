package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/go-go-gadget/pkg/words"
	"github.com/spf13/cobra"
)

var wordCount int

var wordsCmd = &cobra.Command{
	Use:   "words",
	Short: "create list of words",
	Run: func(cmd *cobra.Command, args []string) {
		weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}

		res := words.Words(wordCount, weight)
		fmt.Println(res)
	},
}

func init() {
	wordsCmd.Flags().IntVarP(&wordCount, "count", "c", 8, "Word Count")
	rootCmd.AddCommand(wordsCmd)
}