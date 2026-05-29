package ext

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCaptureRoundTrip(t *testing.T) {
	dir := t.TempDir()
	runIn(t, dir, func() {
		_ = os.Mkdir("src", 0o755)
		_ = os.WriteFile(filepath.Join("src", "a.go"), nil, 0o644)
		_ = os.WriteFile(filepath.Join("src", "b.go"), nil, 0o644)

		out := filepath.Join(dir, "captured.yaml")
		report, err := Capture(CaptureOptions{Name: "rt", OutputPath: out, Depth: -1})
		if err != nil {
			t.Fatalf("Capture: %v", err)
		}
		if report.Dirs != 1 {
			t.Errorf("expected 1 dir, got %d", report.Dirs)
		}
		if report.Files != 2 {
			t.Errorf("expected 2 files, got %d", report.Files)
		}
		if report.Patterns != 1 {
			t.Errorf("expected 1 pattern (uniform .go), got %d", report.Patterns)
		}

		data, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "name: rt") {
			t.Errorf("captured YAML missing name, got:\n%s", data)
		}
	})
}

func TestCaptureWithReporterInvokesCallback(t *testing.T) {
	dir := t.TempDir()
	runIn(t, dir, func() {
		_ = os.Mkdir("docs", 0o755)
		_ = os.WriteFile(filepath.Join("docs", "x.md"), nil, 0o644)

		out := filepath.Join(dir, "out.yaml")
		var seen CaptureReport
		called := false
		_, err := CaptureWithReporter(CaptureOptions{Name: "cb", OutputPath: out}, func(r CaptureReport) {
			called = true
			seen = r
		})
		if err != nil {
			t.Fatal(err)
		}
		if !called {
			t.Fatal("reporter callback not invoked")
		}
		if seen.Dirs != 1 || seen.Files != 1 || seen.Patterns != 1 {
			t.Errorf("unexpected report: %+v", seen)
		}
	})
}

func TestCaptureDepthLimit(t *testing.T) {
	dir := t.TempDir()
	runIn(t, dir, func() {
		_ = os.MkdirAll(filepath.Join("a", "b", "c"), 0o755)
		_ = os.WriteFile(filepath.Join("a", "b", "c", "deep.md"), nil, 0o644)
		_ = os.WriteFile("top.md", nil, 0o644)

		out := filepath.Join(dir, "out.yaml")
		report, err := Capture(CaptureOptions{Name: "limit", OutputPath: out, Depth: 1})
		if err != nil {
			t.Fatal(err)
		}
		if report.Dirs != 1 {
			t.Errorf("expected depth-1 to capture only one dir, got %d", report.Dirs)
		}
		if report.Files != 1 {
			t.Errorf("expected depth-1 to capture only top.md, got %d", report.Files)
		}
	})
}
