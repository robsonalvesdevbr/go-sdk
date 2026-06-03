package sdk

import (
	"os/exec"
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
