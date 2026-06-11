package sdk

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

// writeTarGz builds a gzipped tar archive at path containing the given entries
// (name -> contents). A name ending in "/" is written as a directory entry.
func writeTarGz(t *testing.T, path string, entries map[string]string) {
	t.Helper()

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create archive: %v", err)
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)

	for name, body := range entries {
		hdr := &tar.Header{Name: name, Mode: 0o644}
		if name[len(name)-1] == '/' {
			hdr.Typeflag = tar.TypeDir
			hdr.Mode = 0o755
			if err := tw.WriteHeader(hdr); err != nil {
				t.Fatalf("write dir header: %v", err)
			}
			continue
		}
		hdr.Typeflag = tar.TypeReg
		hdr.Size = int64(len(body))
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("write file header: %v", err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatalf("write file body: %v", err)
		}
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
}

func TestExtractTarGzRejectsPathTraversal(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "evil.tar.gz")
	dst := filepath.Join(tmp, "dst")

	writeTarGz(t, archive, map[string]string{
		"../escape.txt": "pwned",
	})

	if err := ExtractTarGz(archive, dst); err == nil {
		t.Fatal("expected ExtractTarGz to reject path traversal entry, got nil error")
	}

	escaped := filepath.Join(tmp, "escape.txt")
	if _, err := os.Stat(escaped); err == nil {
		t.Fatalf("path traversal succeeded: %s was created", escaped)
	}
}

func TestExtractTarGzExtractsValidArchive(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "go.tar.gz")
	dst := filepath.Join(tmp, "dst")

	writeTarGz(t, archive, map[string]string{
		"go/":        "",
		"go/bin/go":  "binary",
		"go/VERSION": "go1.99.0",
	})

	if err := ExtractTarGz(archive, dst); err != nil {
		t.Fatalf("ExtractTarGz failed on valid archive: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(dst, "go", "VERSION"))
	if err != nil {
		t.Fatalf("expected extracted VERSION file: %v", err)
	}
	if string(got) != "go1.99.0" {
		t.Fatalf("unexpected VERSION contents: %q", got)
	}
}
