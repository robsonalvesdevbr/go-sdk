/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package build

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
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

// defaultMinorCount limits the default output to the newest minor release families (e.g. go1.26.x).
const defaultMinorCount = 5

const (
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

// isTerminal reports whether w is an interactive terminal, so ANSI colors are
// only emitted when a human is looking at the output (pipes stay clean).
func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func newCreateCmdList(versions *[]entity.GoVersion) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List stable Go versions",
		Long:  `Lists stable Go versions available on the official Go repository, newest first. By default only the latest minor releases are shown; use --all to list every version. The current version installed on the system will be marked with "(current)".`,
		RunE:  runCreateList(versions),
	}
	cmd.Flags().BoolP("all", "a", false, "list every stable Go version instead of only the latest minor releases")
	return cmd
}

func runCreateList(versions *[]entity.GoVersion) cli.RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		version, err := sdk.GetSystemGoVersion()
		if err != nil {
			cmd.Println("Error:", err)
			return err
		}

		re := regexp.MustCompile(`go\d+\.\d+(?:\.\d+)?`)
		version = re.FindString(version)
		currentVersion := "v" + strings.TrimPrefix(version, "go")

		showAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		sorted := make([]entity.GoVersion, len(*versions))
		copy(sorted, *versions)
		sort.Slice(sorted, func(i, j int) bool {
			vi := "v" + strings.TrimPrefix(sorted[i].Version, "go")
			vj := "v" + strings.TrimPrefix(sorted[j].Version, "go")
			return semver.Compare(vi, vj) > 0
		})

		cells := make([]string, 0, len(sorted))
		width := 0
		currentIdx := -1
		minors := []string{}
		for _, v := range sorted {
			versionOfList := "v" + strings.TrimPrefix(v.Version, "go")
			if !showAll {
				minor := semver.MajorMinor(versionOfList)
				if !slices.Contains(minors, minor) {
					// sorted is descending, so every version after this one is older
					if len(minors) == defaultMinorCount {
						break
					}
					minors = append(minors, minor)
				}
			}

			cell := v.Version
			if semver.Compare(currentVersion, versionOfList) == 0 {
				cell = fmt.Sprintf("%s (current)", v.Version)
				currentIdx = len(cells)
			}
			if len(cell) > width {
				width = len(cell)
			}
			cells = append(cells, cell)
		}

		out := cmd.OutOrStdout()
		useColor := isTerminal(out)
		const columns = 5
		for i := 0; i < len(cells); i += columns {
			end := min(i+columns, len(cells))
			var sb strings.Builder
			for j, cell := range cells[i:end] {
				if useColor && i+j == currentIdx {
					// pad manually: ANSI escapes must not count toward column width
					padding := strings.Repeat(" ", max(width+2-len(cell), 0))
					sb.WriteString(colorGreen + cell + colorReset + padding)
					continue
				}
				fmt.Fprintf(&sb, "%-*s", width+2, cell)
			}
			fmt.Fprintln(out, strings.TrimRight(sb.String(), " "))
		}

		if !showAll {
			// hint goes to stderr so piped output stays clean
			cmd.Printf("\nShowing the latest %d minor releases. Use --all to list every version.\n", defaultMinorCount)
		}
		return nil
	}
}
