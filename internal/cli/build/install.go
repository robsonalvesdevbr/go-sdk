/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Latest version of Go",
	Long:  `Install Latest version of Go. This command will download and install the latest version of Go on your system.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	installCmd.Flags().BoolP("latest", "l", false, "Install latest version of Go")
	installCmd.Flags().StringP("version-number", "v", "", "Version number to install (e.g. 1.16.3)")
	installCmd.MarkFlagsMutuallyExclusive("latest", "version-number")
}

func NewCommandInstall() *cobra.Command {
	return installCmd
}
