/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"
	"strings"

	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List stable Go versions",
	Long:  `Lists all stable Go versions available on the official Go repository. The current version installed on the system will be marked with "(current)".`,
	Run: func(cmd *cobra.Command, args []string) {
		versions, err := sdk.GetListOfGoVersions()
		if err != nil {
			cmd.Println("Error:", err)
			return
		}
		for _, v := range versions {
			version, err := sdk.GetSystemGoVersion()
			if err != nil {
				cmd.Println("Error:", err)
				return
			}

			if strings.Contains(version, v) {
				v = fmt.Sprintf("%s (current)", v)
			}

			cmd.Println(v)
		}
	},
}

func init() {
}

func NewCommandList() *cobra.Command {
	return listCmd
}
