package ext

import (
	"os"
	"path/filepath"
	"testing"
)

// TestContributingDocsPresent guards the docs listed in
// todo/03-contributing-guidelines.md from silent deletion. It only
// checks presence and a non-trivial size — content correctness is a
// human review job.
func TestContributingDocsPresent(t *testing.T) {
	root := findRepoRoot(t)
	required := []string{
		"CONTRIBUTING.md",
		filepath.Join(".github", "PULL_REQUEST_TEMPLATE.md"),
		filepath.Join(".github", "ISSUE_TEMPLATE", "bug_report.md"),
		filepath.Join(".github", "ISSUE_TEMPLATE", "feature_request.md"),
	}
	for _, rel := range required {
		t.Run(rel, func(t *testing.T) {
			info, err := os.Stat(filepath.Join(root, rel))
			if err != nil {
				t.Fatalf("missing required doc: %v", err)
			}
			if info.Size() < 200 {
				t.Errorf("%s is suspiciously short (%d bytes)", rel, info.Size())
			}
		})
	}
}
