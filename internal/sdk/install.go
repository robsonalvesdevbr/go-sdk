package sdk

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DefaultInstallDir is the destination used when no override is provided.
const DefaultInstallDir = "/usr/local"

// ResolveInstallDir picks the install destination, preferring an explicit
// override, then the GO_SDK_INSTALL_DIR environment variable, then the default.
func ResolveInstallDir(override string) string {
	if override != "" {
		return override
	}
	if dir := strings.TrimSpace(os.Getenv("GO_SDK_INSTALL_DIR")); dir != "" {
		return dir
	}
	return DefaultInstallDir
}

// ensureWritable verifies we can actually write into dir before downloading
// anything, by creating and removing a temporary file. The returned error wraps
// fs.ErrPermission when access is denied, so callers can detect it with errors.Is.
func ensureWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("cannot create %s: %w", dir, err)
	}

	f, err := os.CreateTemp(dir, ".go-sdk-*")
	if err != nil {
		return fmt.Errorf("cannot write to %s: %w", dir, err)
	}
	f.Close()
	return os.Remove(f.Name())
}

func LatestVersion() (string, error) {
	resp, err := http.Get("https://go.dev/VERSION?m=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	return strings.TrimSpace(lines[0]), nil
}

func InstallVersion(version, dir string) error {
	if version == "" {
		var err error
		version, err = LatestVersion()
		if err != nil {
			return err
		}
	}

	dir = ResolveInstallDir(dir)

	// Fail fast (before the ~150MB download) if we cannot write to the target.
	if err := ensureWritable(dir); err != nil {
		return err
	}

	url := "https://go.dev/dl/" + version + ".linux-amd64.tar.gz"
	filename := version + ".tar.gz"

	if err := DownloadFile(filename, url); err != nil {
		return err
	}

	if err := ExtractTarGz(filename, dir); err != nil {
		return err
	}

	return os.Remove(filename)
}

func DownloadFile(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func ExtractTarGz(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}

			f, err := os.OpenFile(
				target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				os.FileMode(hdr.Mode),
			)
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}

			f.Close()
		}
	}
}
