package sdk

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// GetSystemGoVersion executes 'go version' in the host OS and returns the output.
func GetSystemGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Clean up the output by removing the newline character at the end
	cleanOutput := strings.TrimSpace(string(output))
	return cleanOutput, nil
}

// GetUseLocalGoVersion is a placeholder function that simulates retrieving the local Go version.
func GetUseLocalGoVersion() (string, error) {
	cmd := exec.Command("which", "go")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Clean up the output by removing the newline character at the end
	cleanOutput := strings.TrimSpace(string(output))
	return cleanOutput, nil
}

// GetListOfGoVersions is a placeholder function that simulates retrieving a list of available Go versions.
func GetListOfGoVersions() ([]string, error) {
	// 1. Equivalent to: git ls-remote https://github.com/golang/go
	cmd := exec.Command("git", "ls-remote", "--tags", "--refs", "https://github.com/golang/go", "main", "go*")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 2. Equivalent to: egrep, awk, and sed
	// Matches lines ending with refs/tags/goX.Y.Z and captures the goX.Y.Z part
	tagRegex := regexp.MustCompile(`refs/tags/(go\d+\.\d+\.\d+)$`)
	var versions []string

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		matches := tagRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			versions = append(versions, matches[1])
		}
	}

	// 3. Equivalent to: sort -V
	sort.Slice(versions, func(i, j int) bool {
		v1 := parseVersion(versions[i])
		v2 := parseVersion(versions[j])

		// Compare each numeric segment of the version
		for k := 0; k < len(v1) && k < len(v2); k++ {
			if v1[k] != v2[k] {
				return v1[k] < v2[k]
			}
		}

		// If all segments match so far, the shorter version comes first
		return len(v1) < len(v2)
	})

	return versions, nil
}

// parseVersion converts a version string like "go1.20.1" into a slice of integers [1, 20, 1].
func parseVersion(version string) []int {
	version = strings.TrimPrefix(version, "go")
	parts := strings.Split(version, ".")

	nums := make([]int, len(parts))
	for i, p := range parts {
		nums[i], _ = strconv.Atoi(p)
	}

	return nums
}
