package sdk

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching latest version: unexpected status %s", resp.Status)
	}

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

	filename := version + ".linux-amd64.tar.gz"
	url := "https://go.dev/dl/" + filename

	// expectedSha256 is the checksum published by go.dev for this archive. It may
	// be empty when the caller could not resolve it; in that case integrity is
	// not verified and DownloadFile only checks transport-level success.
	expectedSha256 := lookupSha256(version, filename)

	tmpFile, err := DownloadFile(url, expectedSha256)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	return ExtractTarGz(tmpFile, dir)
}

// lookupSha256 returns the SHA256 published by go.dev for the given archive
// filename, or "" if it cannot be resolved.
func lookupSha256(version, filename string) string {
	versions, err := GetListOfGoVersionsV2()
	if err != nil {
		return ""
	}
	for _, v := range versions {
		if v.Version != version {
			continue
		}
		for _, f := range v.Files {
			if f.Filename == filename {
				return f.Sha256
			}
		}
	}
	return ""
}

// DownloadFile streams url into a securely-created temporary file and returns
// its path. When expectedSha256 is non-empty, the download is verified against
// it and the temp file is removed on mismatch. The caller owns the returned
// file and is responsible for removing it.
func DownloadFile(url, expectedSha256 string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("downloading %s: unexpected status %s", url, resp.Status)
	}

	file, err := os.CreateTemp("", "go-sdk-*.tar.gz")
	if err != nil {
		return "", err
	}
	tmpPath := file.Name()

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(file, hasher), resp.Body); err != nil {
		file.Close()
		os.Remove(tmpPath)
		return "", err
	}
	if err := file.Close(); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	if expectedSha256 != "" {
		got := hex.EncodeToString(hasher.Sum(nil))
		if !strings.EqualFold(got, expectedSha256) {
			os.Remove(tmpPath)
			return "", fmt.Errorf("checksum mismatch for %s: expected %s, got %s", url, expectedSha256, got)
		}
	}

	return tmpPath, nil
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

	cleanDst := filepath.Clean(dst)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		// Guard against path traversal ("Zip Slip"): reject any entry whose
		// resolved path escapes the destination directory.
		target := filepath.Join(cleanDst, hdr.Name)
		if target != cleanDst && !strings.HasPrefix(target, cleanDst+string(os.PathSeparator)) {
			return fmt.Errorf("refusing to extract %q: path escapes %s", hdr.Name, cleanDst)
		}

		mode := os.FileMode(hdr.Mode).Perm()

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, mode); err != nil {
				return err
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}

			f, err := os.OpenFile(
				target,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				mode,
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
