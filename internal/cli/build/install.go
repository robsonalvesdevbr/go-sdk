/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"
	"slices"

	"github.com/robsonalvesdevbr/go-sdk/internal/cli"
	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"github.com/spf13/cobra"
)

var installCmd *cobra.Command

func NewCommandInstall(versions *[]string) *cobra.Command {
	installCmd = newCreateCmdInstall(versions)
	installCmd.Flags().BoolP("latest", "l", false, "Install latest version of Go")
	installCmd.Flags().StringP("version-number", "v", "", "Version number to install (e.g. 1.16.3)")
	installCmd.MarkFlagsMutuallyExclusive("latest", "version-number")
	return installCmd
}

func newCreateCmdInstall(versions *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install a specific version of Go",
		Long:  `Install a specific version of Go. This command will download and install the specified version of Go on your system.`,
		RunE:  runCreateInstall(versions),
	}
	return cmd
}

func runCreateInstall(versions *[]string) cli.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		latest, _ := cmd.Flags().GetBool("latest")
		versionNumber := cmd.Flag("version-number").Value.String()

		if !slices.Contains(*versions, fmt.Sprintf("go%s", versionNumber)) {
			return fmt.Errorf("version %s is not available", versionNumber)
		}

		if latest {
			fmt.Println("Installing latest version of Go...")
			err := sdk.InstallVersion("")
			if err != nil {
				fmt.Printf("Error installing Go: %v\n", err)
				return err
			} else {
				fmt.Println("Go installed successfully!")
			}
		} else if versionNumber != "" {
			fmt.Printf("Installing Go version %s...\n", versionNumber)
			err := sdk.InstallVersion(fmt.Sprintf("go%s", versionNumber))
			if err != nil {
				fmt.Printf("Error installing Go: %v\n", err)
				return err
			} else {
				fmt.Println("Go installed successfully!")
			}
		} else {
			fmt.Println("Please specify either --latest or --version-number flag.")
			return fmt.Errorf("please specify either --latest or --version-number flag")
		}
		return nil
	}
}
