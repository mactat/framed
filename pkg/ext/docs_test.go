package ext

import (
	"os"
	"path/filepath"
	"strings"
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

// TestReadmeAcceptanceCriteria pins the README acceptance criteria
// from todo/07-improve-readme.md so a future rewrite can't silently
// regress them.
func TestReadmeAcceptanceCriteria(t *testing.T) {
	root := findRepoRoot(t)
	body, err := os.ReadFile(filepath.Join(root, "README.md"))
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	text := string(body)

	// "TODO checklist is gone from the README; replaced by a one-line
	// pointer to todo/."
	for _, gone := range []string{"## TODO", "- [ ] ", "- [x] "} {
		if strings.Contains(text, gone) {
			t.Errorf("README still carries the old TODO checklist marker: %q", gone)
		}
	}
	if !strings.Contains(text, "todo/") {
		t.Errorf("README should link to todo/")
	}

	// "No placeholder links."
	if strings.Contains(text, "link-to-full-docs") {
		t.Error("README still references the link-to-full-docs placeholder")
	}

	// "README is meaningfully shorter than today (target: under 220 lines)."
	lines := strings.Count(text, "\n") + 1
	if lines > 220 {
		t.Errorf("README is %d lines, target is under 220", lines)
	}
}
