package go_go_gadget

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/nretnilkram/go-go-gadget/pkg/m8s"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
)

type m8sImageOptions struct {
	useAlpine      bool
	useBusyBox     bool
	useUbuntu      bool
	namespace      string
	containerImage string
}

func resolveImage(o m8sImageOptions) (image, pkgManager string) {
	switch {
	case o.containerImage != "":
		return o.containerImage, "unknown"
	case o.useBusyBox:
		return "busybox", "none"
	case o.useUbuntu:
		return "ubuntu", "apt"
	default:
		// Alpine is the default; --alpine is accepted for CLI compatibility.
		_ = o.useAlpine
		return "alpine", "apk"
	}
}

func bindM8sImageFlags(cmd *cobra.Command, o *m8sImageOptions) {
	cmd.Flags().BoolVarP(&o.useAlpine, "alpine", "a", false, "Use Alpine")
	cmd.Flags().BoolVarP(&o.useBusyBox, "busybox", "b", false, "Use BusyBox")
	cmd.Flags().BoolVarP(&o.useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "default", "Namespace to deploy to")
	cmd.Flags().StringVarP(&o.containerImage, "image", "i", "", "Image to use")
}

func newM8sCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "m8s",
		Short: "Useful k8s commands",
		Long:  `Set of useful kubectl utilities.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newM8sDeploymentCmd())
	cmd.AddCommand(newM8sPodCmd())
	cmd.AddCommand(newM8sConnectionCmd())
	return cmd
}

func newM8sDeploymentCmd() *cobra.Command {
	var opts m8sImageOptions

	cmd := &cobra.Command{
		Use:     "deployment",
		Aliases: []string{"d"},
		Short:   "Create Deployment",
		Long: `Create a Deployment in current kubectl context

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: deployment, d`,
		RunE: func(cmd *cobra.Command, args []string) error {
			image, pkgManager := resolveImage(opts)
			resourceName := "m8-" + m8s.Image2Name(image, "-")
			m8s.PrintInfo(cmd.OutOrStdout(), "deployment", image, pkgManager, resourceName)
			output, err := utilities.RunCommandWithError("kubectl", []string{"create", "deployment", resourceName, "--namespace", opts.namespace, "--image", image, "--", "tail", "-f", "/dev/null"}, ".")
			write(cmd.OutOrStdout(), output)
			return err
		},
	}
	bindM8sImageFlags(cmd, &opts)
	return cmd
}

func newM8sPodCmd() *cobra.Command {
	var opts m8sImageOptions

	cmd := &cobra.Command{
		Use:     "pod",
		Aliases: []string{"p"},
		Short:   "Create Pod",
		Long: `Create a Pod in current kubectl context

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: pod, p`,
		RunE: func(cmd *cobra.Command, args []string) error {
			image, pkgManager := resolveImage(opts)
			resourceName := "m8-utility-" + m8s.Image2Name(image, "-")
			m8s.PrintInfo(cmd.OutOrStdout(), "pod", image, pkgManager, resourceName)
			output, err := utilities.RunCommandWithError("kubectl", []string{"run", resourceName, "--namespace", opts.namespace, "--image", image, "--", "tail", "-f", "/dev/null"}, ".")
			write(cmd.OutOrStdout(), output)
			return err
		},
	}
	bindM8sImageFlags(cmd, &opts)
	return cmd
}

func newM8sConnectionCmd() *cobra.Command {
	var (
		opts             m8sImageOptions
		containerCommand string
	)

	cmd := &cobra.Command{
		Use:     "connection",
		Aliases: []string{"c", "terminal", "t"},
		Short:   "Create Connection",
		Long: `Create a Connection in current kubectl context

Image selection (default: alpine):
  --alpine|-a, --busybox|-b, --ubuntu|-u
  --image|-i <image>   e.g. -i debian

Command (optional):
  --command|-c <command>   e.g. -c /bin/bash
  m8s connection -i public.ecr.aws/aws-cli/aws-cli -c "aws --version"
  m8s connection -i public.ecr.aws/aws-cli/aws-cli -c "/bin/sh"

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: connection, c, terminal, t`,
		RunE: func(cmd *cobra.Command, args []string) error {
			image, pkgManager := resolveImage(opts)
			resourceName := "m8-tmp-utility-" + m8s.Image2Name(image, "-")
			m8s.PrintInfo(cmd.OutOrStdout(), "connection", image, pkgManager, resourceName)
			kubectlArgs := []string{"run", "-it", "--rm", resourceName, "--namespace", opts.namespace, "--image", image}
			if command := strings.Fields(containerCommand); len(command) > 0 {
				kubectlArgs = append(kubectlArgs, append([]string{"--command", "--"}, command...)...)
			}
			return utilities.RunCommandInteractive("kubectl", kubectlArgs, ".")
		},
	}
	bindM8sImageFlags(cmd, &opts)
	cmd.Flags().StringVarP(&containerCommand, "command", "c", "", "Command to run in the container")
	return cmd
}
