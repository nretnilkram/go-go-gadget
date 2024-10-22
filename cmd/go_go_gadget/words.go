package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/go-go-gadget/pkg/words"
	"github.com/spf13/cobra"
)

var wordCount int
var adjectiveWeight int
var animalWeight int
var colorWeight int
var nounWeight int
var verbWeight int

var wordsCmd = &cobra.Command{
	Use:     "words",
	Aliases: []string{"w"},
	Short:   "Create list of words",
	Long: `Generate a list of random words including Adjectives, Animals, Colors, Nouns and Verbs.

Word types are evenly distributed unless a weight is specified. If a weigh is passed in then 0 is used for the weight of missing word types.
Order is always random.

example:
	go-go-gadget words --adjective-weight 1 --animal-weight 3 --verb-weight 1

Aliases: words, w`,
	Run: func(cmd *cobra.Command, args []string) {
		if adjectiveWeight == 0 && animalWeight == 0 && colorWeight == 0 && nounWeight == 0 && verbWeight == 0 {
			weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
			fmt.Println(words.Words(wordCount, weight))
		} else {
			weight := words.WordSetWeight{Adjectives: adjectiveWeight, Animals: animalWeight, Colors: colorWeight, Nouns: nounWeight, Verbs: verbWeight}
			fmt.Println(words.Words(wordCount, weight))
		}
	},
}

func init() {
	wordsCmd.Flags().IntVarP(&wordCount, "count", "c", 8, "Word Count")
	wordsCmd.Flags().IntVarP(&adjectiveWeight, "adjective-weight", "", 0, "Adjective Weight")
	wordsCmd.Flags().IntVarP(&animalWeight, "animal-weight", "", 0, "Animal Weight")
	wordsCmd.Flags().IntVarP(&colorWeight, "color-weight", "", 0, "Color Weight")
	wordsCmd.Flags().IntVarP(&nounWeight, "noun-weight", "", 0, "Noun Weight")
	wordsCmd.Flags().IntVarP(&verbWeight, "verb-weight", "", 0, "Verb Weight")
	rootCmd.AddCommand(wordsCmd)
}
