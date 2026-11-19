package go_go_gadget

import (
	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/words"
)

func newWordsCmd() *cobra.Command {
	var (
		wordCount       int
		adjectiveWeight int
		animalWeight    int
		colorWeight     int
		nounWeight      int
		verbWeight      int
	)

	cmd := &cobra.Command{
		Use:     "words",
		Aliases: []string{"w"},
		Short:   "Create list of words",
		Long: `Generate a list of random words including Adjectives, Animals, Colors, Nouns and Verbs.

Word types are evenly distributed unless a weight is specified. If a weigh is passed in then 0 is used for the weight of missing word types.
Order is always random.

example:
	go-go-gadget words --adjective-weight 1 --animal-weight 3 --verb-weight 1

Aliases: words, w`,
		RunE: func(cmd *cobra.Command, args []string) error {
			weight := words.WordSetWeight{
				Adjectives: adjectiveWeight,
				Animals:    animalWeight,
				Colors:     colorWeight,
				Nouns:      nounWeight,
				Verbs:      verbWeight,
			}
			// If no types are passed in then use an equal weight across all types
			if adjectiveWeight == 0 && animalWeight == 0 && colorWeight == 0 && nounWeight == 0 && verbWeight == 0 {
				weight = words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
			}
			writeln(cmd.OutOrStdout(), words.Words(wordCount, weight))
			return nil
		},
	}

	cmd.Flags().IntVarP(&wordCount, "count", "c", 8, "Word Count")
	cmd.Flags().IntVarP(&adjectiveWeight, "adjective-weight", "", 0, "Adjective Weight")
	cmd.Flags().IntVarP(&animalWeight, "animal-weight", "", 0, "Animal Weight")
	cmd.Flags().IntVarP(&colorWeight, "color-weight", "", 0, "Color Weight")
	cmd.Flags().IntVarP(&nounWeight, "noun-weight", "", 0, "Noun Weight")
	cmd.Flags().IntVarP(&verbWeight, "verb-weight", "", 0, "Verb Weight")
	return cmd
}
