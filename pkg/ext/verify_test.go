package ext

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTemplate writes a template file in a temp dir and returns the
// dir + template path. The temp dir is the layout under which Verify
// is expected to run, because templates use paths relative to the
// CWD.
func writeTemplate(t *testing.T, body string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	tmpl := filepath.Join(dir, "framed.yaml")
	if err := os.WriteFile(tmpl, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir, tmpl
}

func runIn(t *testing.T, dir string, fn func()) {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(cwd) }()
	fn()
}

func TestVerifyOK(t *testing.T) {
	tmpl := `name: ok
structure:
  name: root
  files:
    - README.md
  dirs:
    - name: child
      files:
        - inner.md
`
	dir, path := writeTemplate(t, tmpl)
	if err := os.WriteFile(filepath.Join(dir, "README.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "child"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "child", "inner.md"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	runIn(t, dir, func() {
		report, err := Verify(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !report.OK {
			t.Fatalf("expected OK, got failures: %+v", report.Failures)
		}
	})
}

func TestVerifyMissingFileAndDir(t *testing.T) {
	tmpl := `name: missing
structure:
  name: root
  files:
    - README.md
  dirs:
    - name: child
`
	dir, path := writeTemplate(t, tmpl)
	// neither README.md nor child/ exist

	runIn(t, dir, func() {
		report, err := Verify(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if report.OK {
			t.Fatal("expected failures")
		}
		joined := flatten(report.Failures)
		if !strings.Contains(joined, "File not found") || !strings.Contains(joined, "README.md") {
			t.Errorf("missing file failure not surfaced: %s", joined)
		}
		if !strings.Contains(joined, "Directory not found") || !strings.Contains(joined, "child") {
			t.Errorf("missing dir failure not surfaced: %s", joined)
		}
	})
}

func TestVerifyMinMaxCount(t *testing.T) {
	tmpl := `name: counts
structure:
  name: root
  dirs:
    - name: tight
      minCount: 2
      maxCount: 2
`
	dir, path := writeTemplate(t, tmpl)
	tight := filepath.Join(dir, "tight")
	if err := os.Mkdir(tight, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tight, "a"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected min count failure")
		}
		if !strings.Contains(flatten(report.Failures), "Min count") {
			t.Errorf("expected min count failure, got %+v", report.Failures)
		}
	})

	// Now add two more files -> exceeds max count
	if err := os.WriteFile(filepath.Join(tight, "b"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tight, "c"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected max count failure")
		}
		if !strings.Contains(flatten(report.Failures), "Max count") {
			t.Errorf("expected max count failure, got %+v", report.Failures)
		}
	})
}

func TestVerifyAllowedPatterns(t *testing.T) {
	tmpl := `name: patterns
structure:
  name: root
  dirs:
    - name: only-md
      allowedPatterns:
        - "\\.md$"
`
	dir, path := writeTemplate(t, tmpl)
	target := filepath.Join(dir, "only-md")
	_ = os.Mkdir(target, 0o755)
	_ = os.WriteFile(filepath.Join(target, "ok.md"), nil, 0o644)
	_ = os.WriteFile(filepath.Join(target, "stray.txt"), nil, 0o644)

	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected allowed-pattern failure")
		}
		if !strings.Contains(flatten(report.Failures), "Not all files match required pattern") {
			t.Errorf("expected pattern failure, got %+v", report.Failures)
		}
	})
}

func TestVerifyForbiddenPatterns(t *testing.T) {
	tmpl := `name: forbidden
structure:
  name: root
  dirs:
    - name: no-tmp
      forbiddenPatterns:
        - "\\.tmp$"
`
	dir, path := writeTemplate(t, tmpl)
	target := filepath.Join(dir, "no-tmp")
	_ = os.Mkdir(target, 0o755)
	_ = os.WriteFile(filepath.Join(target, "leak.tmp"), nil, 0o644)

	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected forbidden-pattern failure")
		}
		if !strings.Contains(flatten(report.Failures), "Forbidden pattern") {
			t.Errorf("expected forbidden failure, got %+v", report.Failures)
		}
	})
}

func TestVerifyAllowChildren(t *testing.T) {
	tmpl := `name: nochildren
structure:
  name: root
  dirs:
    - name: flat
      allowChildren: false
`
	dir, path := writeTemplate(t, tmpl)
	flat := filepath.Join(dir, "flat")
	_ = os.Mkdir(flat, 0o755)
	_ = os.Mkdir(filepath.Join(flat, "nested"), 0o755)

	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected allowChildren failure")
		}
		if !strings.Contains(flatten(report.Failures), "Children not allowed") {
			t.Errorf("expected children-not-allowed failure, got %+v", report.Failures)
		}
	})
}

func TestVerifyMaxDepth(t *testing.T) {
	tmpl := `name: depth
structure:
  name: root
  dirs:
    - name: shallow
      maxDepth: 1
`
	dir, path := writeTemplate(t, tmpl)
	shallow := filepath.Join(dir, "shallow")
	_ = os.MkdirAll(filepath.Join(shallow, "a", "b"), 0o755)

	runIn(t, dir, func() {
		report, _ := Verify(path)
		if report.OK {
			t.Fatal("expected maxDepth failure")
		}
		if !strings.Contains(flatten(report.Failures), "Max depth exceeded") {
			t.Errorf("expected max-depth failure, got %+v", report.Failures)
		}
	})
}

func TestVerifyMissingTemplate(t *testing.T) {
	_, err := Verify(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func flatten(failures []Failure) string {
	var b strings.Builder
	for _, f := range failures {
		b.WriteString(f.Prompt)
		b.WriteString(" ")
		b.WriteString(f.Detail)
		b.WriteString("\n")
	}
	return b.String()
}
