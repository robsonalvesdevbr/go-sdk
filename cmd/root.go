/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/robsonalvesdevbr/go-sdk/internal/cli/build"
	"github.com/robsonalvesdevbr/go-sdk/internal/sdk"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-sdk",
	Short: "Manage Golang SDK",
	Long:  `Manage Golang SDK`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var ListVersion []string

func init() {
	ListVersion, _ = sdk.GetListOfGoVersions()
	rootCmd.AddCommand(build.NewCommandCurrent())
	rootCmd.AddCommand(build.NewCommandList(&ListVersion))
	rootCmd.AddCommand(build.NewCommandInstall(&ListVersion))
}
