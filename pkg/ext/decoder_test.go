package ext

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadConfigEDefaults(t *testing.T) {
	tmpl := `name: defaults
structure:
  name: root
`
	dir, path := writeTemplate(t, tmpl)
	_ = dir

	cfg, dirs, err := ReadConfigE(path)
	if err != nil {
		t.Fatalf("ReadConfigE: %v", err)
	}
	if cfg.Name != "defaults" {
		t.Errorf("wrong name: %s", cfg.Name)
	}
	if len(dirs) != 1 {
		t.Fatalf("expected one dir, got %d", len(dirs))
	}
	root := dirs[0]
	if !root.AllowChildren {
		t.Error("AllowChildren should default to true")
	}
	if root.MaxCount == 0 {
		t.Error("MaxCount default should be non-zero")
	}
	if root.MaxDepth == 0 {
		t.Error("MaxDepth default should be non-zero")
	}
	if root.Path != "." {
		t.Errorf("root path should be '.', got %q", root.Path)
	}
}

func TestReadConfigEMissingFile(t *testing.T) {
	if _, _, err := ReadConfigE(filepath.Join(t.TempDir(), "nope.yaml")); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadConfigENoStructure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "noStruct.yaml")
	if err := os.WriteFile(path, []byte("name: bare"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadConfigE(path); err == nil {
		t.Fatal("expected error when structure missing")
	}
}

func TestReadConfigEMalformed(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("::: not yaml ::: ["), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadConfigE(path); err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestReadConfigENestedDirs(t *testing.T) {
	tmpl := `name: nested
structure:
  name: root
  dirs:
    - name: a
      dirs:
        - name: b
          dirs:
            - name: c
`
	_, path := writeTemplate(t, tmpl)
	_, dirs, err := ReadConfigE(path)
	if err != nil {
		t.Fatalf("ReadConfigE: %v", err)
	}
	if len(dirs) != 4 {
		t.Fatalf("expected 4 dirs (root + a + b + c), got %d", len(dirs))
	}
	if dirs[3].Path != "./a/b/c" {
		t.Errorf("deepest path wrong: %q", dirs[3].Path)
	}
}
