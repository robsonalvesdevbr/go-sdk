package sdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	// Use exec.LookPath to find the `go` executable in a cross-platform way
	path, err := exec.LookPath("go")
	if err != nil {
		return "", err
	}

	return path, nil
}

// GetListOfGoVersions is a placeholder function that simulates retrieving a list of available Go versions.
func GetListOfGoVersions() ([]string, error) {
	// Use GitHub API to list tags so we don't require `git` on the host.
	// Optional: set GITHUB_TOKEN env var to increase rate limits.
	client := &http.Client{}
	perPage := 100
	page := 1
	tagRegex := regexp.MustCompile(`^go\d+\.\d+\.\d+$`)
	versions := []string{}
	seen := map[string]bool{}

	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))

	for {
		url := fmt.Sprintf("https://api.github.com/repos/golang/go/tags?per_page=%d&page=%d", perPage, page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		if token != "" {
			req.Header.Set("Authorization", "token "+token)
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("github API returned %d: %s", resp.StatusCode, string(body))
		}

		var items []struct {
			Name string `json:"name"`
		}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&items); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(items) == 0 {
			break
		}

		for _, it := range items {
			if tagRegex.MatchString(it.Name) && !seen[it.Name] {
				versions = append(versions, it.Name)
				seen[it.Name] = true
			}
		}

		if len(items) < perPage {
			break
		}
		page++
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
