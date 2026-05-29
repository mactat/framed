package ext

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirsOnly(t *testing.T) {
	tmpl := `name: create
structure:
  name: root
  files:
    - top.md
  dirs:
    - name: sub
      files:
        - inner.md
`
	dir, path := writeTemplate(t, tmpl)

	runIn(t, dir, func() {
		report, err := Create(CreateOptions{TemplatePath: path, CreateFiles: false})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if report.DirsCreated == 0 {
			t.Fatal("expected at least one dir")
		}
		if report.FilesCreated != 0 {
			t.Errorf("expected zero files (CreateFiles=false), got %d", report.FilesCreated)
		}
		if !DirExists("sub") {
			t.Error("expected sub/ to exist")
		}
		if FileExists("top.md") || FileExists("sub/inner.md") {
			t.Error("expected no files when CreateFiles=false")
		}
	})
}

func TestCreateWithFiles(t *testing.T) {
	tmpl := `name: create-files
structure:
  name: root
  files:
    - top.md
  dirs:
    - name: sub
      files:
        - inner.md
`
	dir, path := writeTemplate(t, tmpl)

	runIn(t, dir, func() {
		report, err := Create(CreateOptions{TemplatePath: path, CreateFiles: true})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if report.FilesCreated != 2 {
			t.Errorf("expected 2 files, got %d", report.FilesCreated)
		}
		for _, want := range []string{"top.md", filepath.Join("sub", "inner.md")} {
			if !FileExists(want) {
				t.Errorf("expected file %s to exist", want)
			}
		}
	})
}

func TestCreateMissingTemplate(t *testing.T) {
	_, err := Create(CreateOptions{TemplatePath: filepath.Join(t.TempDir(), "nope.yaml")})
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestCreateInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("::: not yaml :::"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Create(CreateOptions{TemplatePath: path}); err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}
