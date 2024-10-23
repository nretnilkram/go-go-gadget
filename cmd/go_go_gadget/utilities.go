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
	Aliases: []string{"u", "util", "utils"},
	Short:   "Useful utility commands",
	Long: `Set of useful cli utilities.

Aliases: utilities, u, util, utils`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var gitQuickBranchCmd = &cobra.Command{
	Use:     "quick-branch",
	Aliases: []string{"qb"},
	Short:   "Create a Branch",
	Long: `Create a git branch using the today's date and user input or random words.

Aliases: quick-branch, qb`,
	Run: func(cmd *cobra.Command, args []string) {
		today := now.ShowDateTime("dash", false)

		// Create a name with a color followed by an animal
		color := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 1, Nouns: 0, Verbs: 0})
		animal := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 1, Colors: 0, Nouns: 0, Verbs: 0})
		name := fmt.Sprintf("%s%s", color, animal)

		if len(args) > 0 {
			name = strings.Replace(strings.TrimSpace(strings.Join(args, " ")), " ", "-", -1)
		}

		branchName := "m8-" + today + "-" + strings.Replace(strings.TrimSpace(name), " ", "-", -1)

		command := "git checkout -b " + branchName
		path := utilities.GetWorkingDir()

		if utilities.IsGitRepo(path) {
			fmt.Print(utilities.RunShellCommand(command, path))
		} else {
			fmt.Print("This is not a git repository.")
		}
	},
}

var gitQuickCommitCmd = &cobra.Command{
	Use:     "quick-commit",
	Aliases: []string{"qc"},
	Short:   "Create a commit message",
	Long: `Commit  and push current changes using a random 5 word commit message.

Aliases: quick-commit, qc`,
	Run: func(cmd *cobra.Command, args []string) {
		weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
		message := strings.TrimSpace(words.Words(5, weight))

		command := fmt.Sprintf(`git commit -am "%s"`, message)
		path := utilities.GetWorkingDir()

		if utilities.IsGitRepo(path) {
			fmt.Print(utilities.RunShellCommand(command, path))
			fmt.Print(utilities.RunShellCommand("git push", path))
		} else {
			fmt.Print("This is not a git repository.")
		}
	},
}

var gitEmptyCommitCmd = &cobra.Command{
	Use:     "empty-commit",
	Aliases: []string{"ec"},
	Short:   "Create an empty commit message",
	Long: `Create and push an empty commit with random 3 word commit message.

Aliases: empty-commit, ec`,
	Run: func(cmd *cobra.Command, args []string) {
		weight := words.WordSetWeight{Adjectives: 1, Animals: 1, Colors: 1, Nouns: 1, Verbs: 1}
		message := strings.TrimSpace(words.Words(3, weight))
		command := fmt.Sprintf(`git commit --allow-empty -m "empty commit %s"`, message)
		path := utilities.GetWorkingDir()

		if utilities.IsGitRepo(path) {
			fmt.Print(utilities.RunShellCommand(command, path))
			fmt.Print(utilities.RunShellCommand("git push", path))
		} else {
			fmt.Print("This is not a git repository.")
		}
	},
}

var tfListResourcesCmd = &cobra.Command{
	Use:     "tf-list-resources",
	Aliases: []string{"tflr"},
	Short:   "List Terraform resources from files",
	Long: `Take a list of files and return a list of resources that can be targeted with terraform command.

--target=module.example_1 \
--target="aws_key_pair.example_2 \
--target="aws_security_group.example_3 \

Aliases: tf-list-resources, tflr`,
	Run: func(cmd *cobra.Command, args []string) {
		utilities.ListTFResources(args)
	},
}

var isSemverCmd = &cobra.Command{
	Use:   "semver",
	Short: "Valid Semvar String",
	Long: `Takes a string and returns true or false on whether the string is valid semver.

e.g. 1.2.0`,
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(utilities.RegexTest(args[0], utilities.SemverRegex))
	},
}

func init() {
	utilCmd.AddCommand(gitQuickBranchCmd)

	utilCmd.AddCommand(gitQuickCommitCmd)

	utilCmd.AddCommand(gitEmptyCommitCmd)

	utilCmd.AddCommand(tfListResourcesCmd)

	utilCmd.AddCommand(isSemverCmd)

	rootCmd.AddCommand(utilCmd)
}
