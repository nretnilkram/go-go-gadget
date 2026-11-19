package go_go_gadget

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/nretnilkram/go-go-gadget/pkg/words"
)

func newUtilCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "utilities",
		Aliases: []string{"u", "util", "utils"},
		Short:   "Useful utility commands",
		Long: `Set of useful cli utilities.

Aliases: utilities, u, util, utils`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newGitQuickBranchCmd())
	cmd.AddCommand(newGitQuickCommitCmd())
	cmd.AddCommand(newGitEmptyCommitCmd())
	cmd.AddCommand(newTFListResourcesCmd())
	cmd.AddCommand(newIsSemverCmd())
	cmd.AddCommand(newBumpChartCmd())
	return cmd
}

func newGitQuickBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "quick-branch",
		Aliases: []string{"qb"},
		Short:   "Create a Branch",
		Long: `Create a git branch using the today's date and user input or random words.

Aliases: quick-branch, qb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			today := utilities.ShowDateTime("dash", false)

			color := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 1, Nouns: 0, Verbs: 0})
			animal := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 1, Colors: 0, Nouns: 0, Verbs: 0})
			name := fmt.Sprintf("%s%s", color, animal)

			if len(args) > 0 {
				name = strings.ReplaceAll(strings.TrimSpace(strings.Join(args, " ")), " ", "-")
			}

			branchName := "m8-" + today + "-" + strings.ReplaceAll(strings.TrimSpace(name), " ", "-")

			path, err := utilities.GetWorkingDir()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if utilities.IsGitRepo(path) {
				output, err := utilities.RunCommandWithError("git", []string{"checkout", "-b", branchName}, path)
				write(out, output)
				return err
			}
			writeln(out, "This is not a git repository.")
			return nil
		},
	}
}

func newGitQuickCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "quick-commit",
		Aliases: []string{"qc"},
		Short:   "Create a commit message",
		Long: `Commit  and push current changes using a random 5 word commit message.

Aliases: quick-commit, qc`,
		RunE: func(cmd *cobra.Command, args []string) error {
			adjective := words.Words(1, words.WordSetWeight{Adjectives: 1, Animals: 0, Colors: 0, Nouns: 0, Verbs: 0})
			animal := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 1, Colors: 0, Nouns: 0, Verbs: 0})
			color := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 1, Nouns: 0, Verbs: 0})
			noun := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 0, Nouns: 1, Verbs: 0})
			verb := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 0, Nouns: 0, Verbs: 1})

			message := strings.TrimSpace(fmt.Sprintf(`%s%s%s%s%s`, adjective, color, animal, verb, noun))

			path, err := utilities.GetWorkingDir()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if utilities.IsGitRepo(path) {
				commitOut, err := utilities.RunCommandWithError("git", []string{"commit", "-am", message}, path)
				write(out, commitOut)
				if err != nil {
					return err
				}
				pushOut, err := utilities.RunCommandWithError("git", []string{"push"}, path)
				write(out, pushOut)
				return err
			}
			writeln(out, "This is not a git repository.")
			return nil
		},
	}
}

func newGitEmptyCommitCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "empty-commit",
		Aliases: []string{"ec"},
		Short:   "Create an empty commit message",
		Long: `Create and push an empty commit with random 3 word commit message.

Aliases: empty-commit, ec`,
		RunE: func(cmd *cobra.Command, args []string) error {
			adjective := words.Words(1, words.WordSetWeight{Adjectives: 1, Animals: 0, Colors: 0, Nouns: 0, Verbs: 0})
			animal := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 1, Colors: 0, Nouns: 0, Verbs: 0})
			color := words.Words(1, words.WordSetWeight{Adjectives: 0, Animals: 0, Colors: 1, Nouns: 0, Verbs: 0})

			message := strings.TrimSpace(fmt.Sprintf(`%s%s%s`, adjective, color, animal))
			path, err := utilities.GetWorkingDir()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if utilities.IsGitRepo(path) {
				commitOut, err := utilities.RunCommandWithError("git", []string{"commit", "--allow-empty", "-m", message}, path)
				write(out, commitOut)
				if err != nil {
					return err
				}
				pushOut, err := utilities.RunCommandWithError("git", []string{"push"}, path)
				write(out, pushOut)
				return err
			}
			writeln(out, "This is not a git repository.")
			return nil
		},
	}
}

func newTFListResourcesCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "tf-list-resources",
		Aliases: []string{"tflr"},
		Short:   "List Terraform resources from files",
		Long: `Take a list of files and return a list of resources that can be targeted with terraform command.

--target=module.example_1 \
--target="aws_key_pair.example_2 \
--target="aws_security_group.example_3 \

Aliases: tf-list-resources, tflr`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return utilities.ListTFResources(cmd.OutOrStdout(), args)
		},
	}
}

func newIsSemverCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "semver",
		Aliases: []string{"sem"},
		Short:   "Valid Semvar String",
		Long: `Takes a string and returns true or false on whether the string is valid semver.

e.g. 1.2.0

Aliases: semver, sem`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			writeln(cmd.OutOrStdout(), utilities.ValidateSemver(args[0]))
			return nil
		},
	}
}

func newBumpChartCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "bump-chart major|minor|patch",
		Aliases: []string{"bc", "helm-bump"},
		Short:   "Bump Helm chart versions",
		Long: `Find Chart.yaml files under the current directory and bump their version
by major, minor, or patch.

Examples:
  go-go-gadget utilities bump-chart patch
  go-go-gadget utilities bump-chart minor
  go-go-gadget utilities bump-chart major

Aliases: bump-chart, bc, helm-bump`,
		Args: func(cmd *cobra.Command, args []string) error {
			printBumpChartUsage := func(reason string) error {
				writef(cmd.ErrOrStderr(), "\n%s\n\nUsage:\n  %s major|minor|patch\n\nExamples:\n  %s patch\n  %s minor\n  %s major\n\n",
					reason, cmd.CommandPath(), cmd.CommandPath(), cmd.CommandPath(), cmd.CommandPath())
				return fmt.Errorf("%s", reason)
			}
			if len(args) != 1 {
				return printBumpChartUsage("missing bump type")
			}
			switch strings.ToLower(args[0]) {
			case "major", "minor", "patch":
				return nil
			default:
				return printBumpChartUsage(fmt.Sprintf("invalid bump type: %s", args[0]))
			}
		},
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := utilities.GetWorkingDir()
			if err != nil {
				return err
			}

			bumpType := strings.ToLower(args[0])
			results, err := utilities.BumpChartsInDir(cwd, bumpType)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			updated := 0
			for _, result := range results {
				if result.Skipped {
					writef(out, "\nSkipping %s (%s)\n", result.File, result.SkipReason)
					continue
				}
				writef(out, "\n%s: %s -> %s (%s)\n", result.File, result.OldVersion, result.NewVersion, bumpType)
				updated++
			}
			writef(out, "\nUpdated %d Chart.yaml file(s)\n\n", updated)
			return nil
		},
	}
}
