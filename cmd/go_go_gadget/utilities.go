package go_go_gadget

import (
	"fmt"
	"strings"

	"github.com/nretnilkram/go-go-gadget/pkg/now"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/nretnilkram/go-go-gadget/pkg/words"
	"github.com/spf13/cobra"
)

var utilCmd = &cobra.Command{
	Use:     "utilities",
	Aliases: []string{"util", "utils"},
	Short:   "Useful Utility Commands",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var gitQuickBranchCmd = &cobra.Command{
	Use:     "quick-branch",
	Aliases: []string{"qb"},
	Short:   "Create a Branch",
	Run: func(cmd *cobra.Command, args []string) {
		weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
		today := now.ShowDateTime("dash", false)
		random := words.Words(2, weight)
		branchName := "m8-" + today + "-" + strings.Replace(strings.TrimSpace(random), " ", "-", -1)

		command := "git checkout -b " + branchName
		path := utilities.GetWorkingDir()

		if utilities.IsGitRepo(path) {
			fmt.Println(utilities.RunShellCommand(command, path))
		} else {
			fmt.Println("This is not a git repository.")
		}
	},
}

var gitQuickCommitCmd = &cobra.Command{
	Use:     "quick-commit",
	Aliases: []string{"qc"},
	Short:   "Create a commit message",
	Run: func(cmd *cobra.Command, args []string) {
		weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
		message := strings.TrimSpace(words.Words(5, weight))

		command := "git commit -am '" + message + "'"
		path := utilities.GetWorkingDir()

		if utilities.IsGitRepo(path) {
			fmt.Println(utilities.RunShellCommand(command, path))
		} else {
			fmt.Println("This is not a git repository.")
		}
	},
}

func init() {
	utilCmd.AddCommand(gitQuickBranchCmd)

	utilCmd.AddCommand(gitQuickCommitCmd)

	rootCmd.AddCommand(utilCmd)
}
