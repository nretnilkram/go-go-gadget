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
		name := words.Words(2, weight)

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
	Run: func(cmd *cobra.Command, args []string) {
		utilities.ListTFResources(args)
	},
}

var isSemverCmd = &cobra.Command{
	Use:                   "semver",
	Short:                 "List Terraform resources from files",
	Args:                  cobra.ExactArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		pattern := "^(?P<major>0|[1-9]\\d*)\\.(?P<minor>0|[1-9]\\d*)\\.(?P<patch>0|[1-9]\\d*)(?:-(?P<prerelease>(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$"

		fmt.Println(utilities.RegexTest(args[0], pattern))
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
