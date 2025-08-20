package OpenApi

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const yamlSpec = `
openapi: 3.0.0
info:
  title: Simple API overview
  version: 2.0.0
paths: {}
`

const jsonSpec = `{
  "openapi": "3.0.0",
  "info": { "title": "Simple API overview", "version": "2.0.0" },
  "paths": {}
}`

const yamlNoTitleSpec = `
openapi: 3.0.0
info:
  version: 1.0.0
paths: {}
`

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed writing temp file %s: %v", name, err)
	}
	return path
}

func TestLoadFile_YAML(t *testing.T) {
	td := t.TempDir()
	path := writeTempFile(t, td, "spec.yaml", yamlSpec)

	got, err := LoadTitle(path)
	if err != nil {
		t.Fatalf("LoadTitle returned an error: %v", err)
	}
	want := "Simple API overview"

	if got != want {
		t.Errorf("LoadTitle returned %q, want %q", got, want)
	}
}

func TestLoadFile_JSON(t *testing.T) {
	td := t.TempDir()
	path := writeTempFile(t, td, "spec.json", jsonSpec)

	got, err := LoadTitle(path)
	if err != nil {
		t.Fatalf("LoadTitle returned an error: %v", err)
	}
	want := "Simple API overview"

	if got != want {
		t.Errorf("LoadTitle returned %q, want %q", got, want)
	}
}

func TestLoadTitle_MissingTitle(t *testing.T) {
	td := t.TempDir()
	path := writeTempFile(t, td, "spec-no-title.yaml", yamlNoTitleSpec)

	_, err := LoadTitle(path)
	if err == nil {
		t.Fatalf("expected error for missing info.title, got nil")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "missing info.title") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPrintTitle_PrintsToStdout(t *testing.T) {
	td := t.TempDir()
	path := writeTempFile(t, td, "spec.yaml", yamlSpec)

	orig := os.Stdout
	defer func() { os.Stdout = orig }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w

	callErr := PrintTitle(path)

	// Close writer to let reader EOF.
	_ = w.Close()
	out, readErr := io.ReadAll(r)
	_ = r.Close()

	if callErr != nil {
		t.Fatalf("PrintTitle returned error: %v", callErr)
	}
	if readErr != nil {
		t.Fatalf("reading stdout: %v", readErr)
	}

	got := strings.TrimSpace(string(out))
	want := "Simple API overview"
	if got != want {
		t.Fatalf("stdout mismatch: got %q, want %q", got, want)
	}
}
