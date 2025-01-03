package go_go_gadget

import (
	"fmt"

	"github.com/nretnilkram/go-go-gadget/pkg/m8s"
	"github.com/nretnilkram/go-go-gadget/pkg/utilities"
	"github.com/spf13/cobra"
)

var useAlpine bool
var useBusyBox bool
var useUbuntu bool
var namespace string
var containerImage string

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

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: deployment, d`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		pkgManager := "apk"
		if containerImage != "" {
			image = containerImage
			pkgManager = "unknown"
		} else if useAlpine {
			image = "alpine"
			pkgManager = "apk"
		} else if useBusyBox {
			image = "busybox"
			pkgManager = "none"
		} else if useUbuntu {
			image = "ubuntu"
			pkgManager = "apt"
		}
		resourceName := "m8-" + m8s.Image2Name(image, "-")
		m8s.PrintInfo("deployment", image, pkgManager, resourceName)
		fmt.Print(utilities.RunShellCommand(fmt.Sprintf("kubectl create deployment %s --namespace %s --image %s -- tail -f /dev/null", resourceName, namespace, image), "."))
	},
}

var m8sPodCmd = &cobra.Command{
	Use:     "pod",
	Aliases: []string{"p"},
	Short:   "Create Pod",
	Long: `Create a Pod in current kubernetes context

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: pod, p`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		pkgManager := "apk"
		if containerImage != "" {
			image = containerImage
			pkgManager = "unknown"
		} else if useAlpine {
			image = "alpine"
			pkgManager = "apk"
		} else if useBusyBox {
			image = "busybox"
			pkgManager = "none"
		} else if useUbuntu {
			image = "ubuntu"
			pkgManager = "apt"
		}
		resourceName := "m8-utility-" + m8s.Image2Name(image, "-")
		m8s.PrintInfo("pod", image, pkgManager, resourceName)
		fmt.Print(utilities.RunShellCommand(fmt.Sprintf("kubectl run %s --namespace %s --image %s -- tail -f /dev/null", resourceName, namespace, image), "."))
	},
}

var m8sConnectionCmd = &cobra.Command{
	Use:     "connection",
	Aliases: []string{"c", "terminal", "t"},
	Short:   "Create Connection",
	Long: `Create a Connection in current kubernetes context

Some Image Options:
* archlinux
* rockylinux
* amazonlinux
* debian

Aliases: connection, c, terminal, t`,
	Run: func(cmd *cobra.Command, args []string) {
		image := "alpine"
		pkgManager := "apk"
		if containerImage != "" {
			image = containerImage
			pkgManager = "unknown"
		} else if useAlpine {
			image = "alpine"
			pkgManager = "apk"
		} else if useBusyBox {
			image = "busybox"
			pkgManager = "none"
		} else if useUbuntu {
			image = "ubuntu"
			pkgManager = "apt"
		}
		resourceName := "m8-tmp-utility-" + m8s.Image2Name(image, "-")
		m8s.PrintInfo("connection", image, pkgManager, resourceName)
		utilities.RunShellCommandInteract(fmt.Sprintf("kubectl run -it --rm %s --namespace %s --image %s", resourceName, namespace, image), ".")
	},
}

func init() {
	m8sDeploymentCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sDeploymentCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sDeploymentCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sDeploymentCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to deploy to")
	m8sDeploymentCmd.Flags().StringVarP(&containerImage, "image", "i", "", "Image to use")
	m8sCmd.AddCommand(m8sDeploymentCmd)

	m8sPodCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sPodCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sPodCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sPodCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to deploy to")
	m8sPodCmd.Flags().StringVarP(&containerImage, "image", "i", "", "Image to use")
	m8sCmd.AddCommand(m8sPodCmd)

	m8sConnectionCmd.Flags().BoolVarP(&useAlpine, "alpine", "a", false, "Use Alpine")
	m8sConnectionCmd.Flags().BoolVarP(&useBusyBox, "busybox", "b", false, "Use BusyBox")
	m8sConnectionCmd.Flags().BoolVarP(&useUbuntu, "ubuntu", "u", false, "Use Ubuntu")
	m8sConnectionCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Namespace to deploy to")
	m8sConnectionCmd.Flags().StringVarP(&containerImage, "image", "i", "", "Image to use")
	m8sCmd.AddCommand(m8sConnectionCmd)

	rootCmd.AddCommand(m8sCmd)
}
