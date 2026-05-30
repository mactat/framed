package ext

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExamplesLoad walks the repo's examples/ directory and asserts
// every example template parses cleanly and round-trips through
// create+verify. Catches drift where someone adds an example that
// can't actually be materialized.
func TestExamplesLoadAndVerify(t *testing.T) {
	repoRoot := findRepoRoot(t)
	examplesDir := filepath.Join(repoRoot, "examples")

	entries, err := os.ReadDir(examplesDir)
	if err != nil {
		t.Fatalf("read examples dir: %v", err)
	}

	any := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		any = true
		example := e.Name()
		t.Run(example, func(t *testing.T) {
			workDir := t.TempDir()
			src := filepath.Join(examplesDir, example)
			body, err := os.ReadFile(src)
			if err != nil {
				t.Fatal(err)
			}
			dst := filepath.Join(workDir, "framed.yaml")
			if err := os.WriteFile(dst, body, 0o644); err != nil {
				t.Fatal(err)
			}

			runIn(t, workDir, func() {
				if _, err := Create(CreateOptions{TemplatePath: dst, CreateFiles: true}); err != nil {
					t.Fatalf("Create: %v", err)
				}
				report, err := Verify(dst)
				if err != nil {
					t.Fatalf("Verify: %v", err)
				}
				if !report.OK {
					t.Fatalf("Verify failed for example %s: %+v", example, report.Failures)
				}
			})
		})
	}
	if !any {
		t.Fatal("no .yaml examples found")
	}
}

// findRepoRoot walks up from the test file's CWD looking for go.mod.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found while walking up from test")
		}
		dir = parent
	}
}
