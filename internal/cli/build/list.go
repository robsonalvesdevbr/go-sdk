/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"
	"strings"

	"github.com/robsonalvesdevbr/go-sdk/internal/cli"
	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd *cobra.Command

func NewCommandList(versions *[]string) *cobra.Command {
	listCmd = newCreateCmdList(versions)
	return listCmd
}

func newCreateCmdList(versions *[]string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stable Go versions",
		Long:  `Lists all stable Go versions available on the official Go repository. The current version installed on the system will be marked with "(current)".`,
		RunE:  runCreateList(versions),
	}
	return cmd
}

func runCreateList(versions *[]string) cli.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		for _, v := range *versions {
			version, err := sdk.GetSystemGoVersion()
			if err != nil {
				cmd.Println("Error:", err)
				return err
			}

			if strings.Contains(version, v) {
				v = fmt.Sprintf("%s (current)", v)
			}

			cmd.Println(v)
		}
		return nil
	}
}
