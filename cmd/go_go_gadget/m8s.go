package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/spf13/cobra"
)

var useAlpine bool
var useBusyBox bool
var useUbuntu bool
var namespace string

var m8sCmd = &cobra.Command{
	Use:   "m8s",
	Short: "Useful k8s commands",
	Long:  `Set of useful kubernetes utilities.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var m8sDeploymentCmd = &cobra.Command{
	Use:     "deployment",
	Aliases: []string{"d"},
	Short:   "Create Deployment",
	Long: `Create a Deployment in current kubernetes context

Aliases: deployment, d`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		if useAlpine {
			image = "alpine"
		} else if useBusyBox {
			image = "busybox"
		} else if useUbuntu {
			image = "ubuntu"
		}
		fmt.Println(utilities.RunShellCommand(fmt.Sprintf("kubectl create deployment m8-%s --namespace %s --image %s -- tail -f /dev/null", image, namespace, image), "."))
	},
}

var m8sPodCmd = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"p"},
	Short:   "Create Pod",
	Long: `Create a Pod in current kubernetes context

Aliases: pod, p`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		if useAlpine {
			image = "alpine"
		} else if useBusyBox {
			image = "busybox"
		} else if useUbuntu {
			image = "ubuntu"
		}
		fmt.Println(utilities.RunShellCommand(fmt.Sprintf("kubectl run m8-utility-%s --namespace %s --image %s -- tail -f /dev/null", image, namespace, image), "."))
	},
}

var m8sConnectionCmd = &cobra.Command{
	Use:     "connection",
	Aliases: []string{"c", "terminal", "t"},
	Short:   "Create Connection",
	Long: `Create a Connection in current kubernetes context

Aliases: connection, c, terminal, t`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		if useAlpine {
			image = "alpine"
		} else if useBusyBox {
			image = "busybox"
		} else if useUbuntu {
			image = "ubuntu"
		}
		utilities.RunShellCommandInteract(fmt.Sprintf("kubectl run -it --rm m8-tmp-utility-%s --namespace %s --image %s", image, namespace, image), ".")
	},
}

func init() {
	m8sDeploymentCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sDeploymentCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sDeploymentCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sDeploymentCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Use Ubuntu")
	m8sCmd.AddCommand(m8sDeploymentCmd)

	m8sPodCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sPodCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sPodCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sPodCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Use Ubuntu")
	m8sCmd.AddCommand(m8sPodCmd)

	m8sConnectionCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sConnectionCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sConnectionCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sConnectionCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Use Ubuntu")
	m8sCmd.AddCommand(m8sConnectionCmd)

	rootCmd.AddCommand(m8sCmd)
}
