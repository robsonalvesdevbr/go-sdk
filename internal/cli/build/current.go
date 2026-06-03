/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"

	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"github.com/spf13/cobra"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Current Go version",
	Long:  `Displays the current Go version installed on the system.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, err := sdk.GetSystemGoVersion()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		cmd.Printf("System Go version: %s\n", version)

		local, _ := cmd.Flags().GetBool("local")
		if local {
			localVersion, err := sdk.GetUseLocalGoVersion()
			if err != nil {
				cmd.Println("Error:", err)
				return
			}
			cmd.Printf("Local Go version: %s\n", localVersion)
		}
	},
}

func init() {
	currentCmd.Flags().BoolP("local", "l", false, "Use local version of Go")
}

func NewCommandCurrent() *cobra.Command {
	return currentCmd
}
