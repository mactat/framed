package ext

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, dir, name, body string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestImportLocalTemplate(t *testing.T) {
	dir := t.TempDir()

	// shared.yaml — the imported subtree.
	write(t, dir, "shared.yaml", `name: shared
structure:
  name: ignored_after_merge
  files:
    - common.md
  dirs:
    - name: nested
      files:
        - inner.md
`)

	// root template references shared.yaml on a child dir.
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: lib
      template: ./shared.yaml
`)

	_, dirs, err := ReadConfigE(rootPath)
	if err != nil {
		t.Fatalf("ReadConfigE: %v", err)
	}

	// Expect: root + lib (host, name preserved) + nested (from import)
	names := dirNames(dirs)
	for _, want := range []string{"root", "lib", "nested"} {
		if !contains(names, want) {
			t.Errorf("expected dir %q in tree, got %v", want, names)
		}
	}
	// The lib dir should have inherited the imported files.
	lib := findByName(dirs, "lib")
	if lib == nil || lib.Files == nil || !contains(*lib.Files, "common.md") {
		t.Errorf("expected lib/ to have imported file common.md, got %+v", lib)
	}
}

func TestImportRemoteTemplate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`name: remote
structure:
  name: r
  files:
    - remote.md
`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: pulled
      template: `+srv.URL+`/remote.yaml
`)

	_, dirs, err := ReadConfigE(rootPath)
	if err != nil {
		t.Fatalf("ReadConfigE: %v", err)
	}
	pulled := findByName(dirs, "pulled")
	if pulled == nil || pulled.Files == nil || !contains(*pulled.Files, "remote.md") {
		t.Errorf("expected pulled/ to have remote.md, got %+v", pulled)
	}
}

func TestImportNestedLocal(t *testing.T) {
	dir := t.TempDir()
	// level2 imported by level1, level1 imported by root.
	write(t, dir, "level2.yaml", `name: level2
structure:
  name: x
  files:
    - leaf.md
`)
	write(t, dir, "level1.yaml", `name: level1
structure:
  name: x
  dirs:
    - name: inner
      template: ./level2.yaml
`)
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: outer
      template: ./level1.yaml
`)

	_, dirs, err := ReadConfigE(rootPath)
	if err != nil {
		t.Fatalf("ReadConfigE: %v", err)
	}
	names := dirNames(dirs)
	for _, want := range []string{"outer", "inner"} {
		if !contains(names, want) {
			t.Errorf("expected dir %q via nested import, got %v", want, names)
		}
	}
	inner := findByName(dirs, "inner")
	if inner == nil || inner.Files == nil || !contains(*inner.Files, "leaf.md") {
		t.Errorf("expected inner/ to have leaf.md from nested import, got %+v", inner)
	}
}

func TestImportCycleDetected(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "a.yaml", `name: a
structure:
  name: x
  dirs:
    - name: hop
      template: ./b.yaml
`)
	write(t, dir, "b.yaml", `name: b
structure:
  name: x
  dirs:
    - name: back
      template: ./a.yaml
`)
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: start
      template: ./a.yaml
`)

	_, _, err := ReadConfigE(rootPath)
	if err == nil {
		t.Fatal("expected cycle error")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected cycle error, got %v", err)
	}
}

func TestImportMissingFile(t *testing.T) {
	dir := t.TempDir()
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: oops
      template: ./does-not-exist.yaml
`)
	_, _, err := ReadConfigE(rootPath)
	if err == nil {
		t.Fatal("expected error for missing imported template")
	}
}

func TestImportBadHTTPStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	dir := t.TempDir()
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: pulled
      template: `+srv.URL+`/missing.yaml
`)
	_, _, err := ReadConfigE(rootPath)
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
	if !strings.Contains(err.Error(), "404") && !strings.Contains(err.Error(), "status") {
		t.Errorf("expected status-code error, got %v", err)
	}
}

func TestImportMalformedYAML(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "broken.yaml", "::: not yaml :::")
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: oops
      template: ./broken.yaml
`)
	_, _, err := ReadConfigE(rootPath)
	if err == nil {
		t.Fatal("expected error for malformed imported YAML")
	}
}

func TestImportRejectsInlineAndTemplate(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "shared.yaml", `name: shared
structure:
  name: x
  files:
    - hello.md
`)
	rootPath := write(t, dir, "framed.yaml", `name: parent
structure:
  name: root
  dirs:
    - name: hybrid
      template: ./shared.yaml
      files:
        - local.md
`)
	_, _, err := ReadConfigE(rootPath)
	if err == nil {
		t.Fatal("expected error when both template and inline files are set")
	}
}

// helpers

func dirNames(dirs []SingleDir) []string {
	out := make([]string, 0, len(dirs))
	for _, d := range dirs {
		out = append(out, d.Name)
	}
	return out
}

func findByName(dirs []SingleDir, name string) *SingleDir {
	for i := range dirs {
		if dirs[i].Name == name {
			return &dirs[i]
		}
	}
	return nil
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}
	return false
}
