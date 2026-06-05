/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"errors"
	"fmt"
	"io/fs"
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
	installCmd.Flags().StringP("dir", "d", "", "Target install directory (default /usr/local or $GO_SDK_INSTALL_DIR)")
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
		dir := sdk.ResolveInstallDir(cmd.Flag("dir").Value.String())

		var version string
		switch {
		case latest:
			fmt.Println("Installing latest version of Go...")
			version = ""
		case versionNumber != "":
			if !slices.Contains(*versions, fmt.Sprintf("go%s", versionNumber)) {
				return fmt.Errorf("version %s is not available", versionNumber)
			}
			fmt.Printf("Installing Go version %s...\n", versionNumber)
			version = fmt.Sprintf("go%s", versionNumber)
		default:
			fmt.Println("Please specify either --latest or --version-number flag.")
			return fmt.Errorf("please specify either --latest or --version-number flag")
		}

		if err := sdk.InstallVersion(version, dir); err != nil {
			if errors.Is(err, fs.ErrPermission) {
				fmt.Printf("Permission denied writing to %s.\n", dir)
				fmt.Println("Re-run with sudo, or pick a writable directory:")
				fmt.Println("  sudo go-sdk install --latest")
				fmt.Println(`  go-sdk install --latest --dir "$HOME/.local"      # then add <dir>/go/bin to PATH`)
				fmt.Println(`  GO_SDK_INSTALL_DIR="$HOME/.local" go-sdk install --latest`)
				return err
			}
			fmt.Printf("Error installing Go: %v\n", err)
			return err
		}

		fmt.Println("Go installed successfully!")
		return nil
	}
}
