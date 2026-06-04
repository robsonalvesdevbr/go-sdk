package cli

import "github.com/spf13/cobra"

type RunEFunc func(cmd *cobra.Command, args []string) error
