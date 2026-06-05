/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/robsonalvesdevbr/go-sdk/internal/cli"
	"github.com/robsonalvesdevbr/go-sdk/internal/entity"
	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"golang.org/x/mod/semver"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd *cobra.Command

func NewCommandList(versions *[]entity.GoVersion) *cobra.Command {
	listCmd = newCreateCmdList(versions)
	return listCmd
}

func newCreateCmdList(versions *[]entity.GoVersion) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stable Go versions",
		Long:  `Lists all stable Go versions available on the official Go repository. The current version installed on the system will be marked with "(current)".`,
		RunE:  runCreateList(versions),
	}
	return cmd
}

func runCreateList(versions *[]entity.GoVersion) cli.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		for _, v := range *versions {
			version, err := sdk.GetSystemGoVersion()
			if err != nil {
				cmd.Println("Error:", err)
				return err
			}

			re := regexp.MustCompile(`go\d+\.\d+(?:\.\d+)?`)
			version = re.FindString(version)

			currentVersion := "v" + strings.TrimPrefix(version, "go")
			versionOfList := "v" + strings.TrimPrefix(v.Version, "go")

			if semver.Compare(currentVersion, versionOfList) == 0 {
				v = entity.GoVersion{Version: fmt.Sprintf("%s (current)", v.Version)}
			}

			cmd.Println(v.Version)
		}
		return nil
	}
}
