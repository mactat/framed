package ext

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// resolveImports walks the structure tree and, for every node that
// declares a `template:` ref, replaces its child subtree with the
// referenced template's structure. Local paths are resolved relative
// to the *importing* template's directory (not the CWD). Remote
// refs (`http://`, `https://`) are fetched with net/http.
//
// Imports nest — a referenced template may itself contain imports.
// Cycles are detected via the `seen` set keyed by absolute path or
// canonical URL.
//
// The merge rule for v1 is "subtree replacement": the imported
// template's structure replaces the host dir's Files/Dirs and the
// scalar fields (counts, patterns, depth, allowChildren). The host
// dir's Name is preserved. Declaring both `template:` and inline
// `dirs:`/`files:` on the same node is rejected so semantics stay
// unambiguous.
func resolveImports(node *SingleDir, basePath string, seen map[string]bool) error {
	if node == nil {
		return nil
	}

	if node.Template != "" {
		if node.Dirs != nil || node.Files != nil {
			return fmt.Errorf("node %q declares both `template:` and inline dirs/files — pick one", node.Name)
		}

		ref, err := canonicalize(node.Template, basePath)
		if err != nil {
			return err
		}
		if seen[ref] {
			return fmt.Errorf("template import cycle detected at %q", ref)
		}
		seen[ref] = true
		defer delete(seen, ref)

		imported, nextBase, err := loadTemplate(node.Template, basePath)
		if err != nil {
			return err
		}
		if imported.Structure == nil {
			return fmt.Errorf("imported template %q has no structure", node.Template)
		}

		if err := resolveImports(imported.Structure, nextBase, seen); err != nil {
			return err
		}

		mergeFromImport(node, imported.Structure)
		node.Template = "" // consumed
	}

	if node.Dirs != nil {
		for i := range *node.Dirs {
			child := &(*node.Dirs)[i]
			if err := resolveImports(child, basePath, seen); err != nil {
				return err
			}
		}
	}
	return nil
}

func mergeFromImport(host *SingleDir, imported *SingleDir) {
	host.Files = imported.Files
	host.Dirs = imported.Dirs
	host.AllowedPatterns = imported.AllowedPatterns
	host.ForbiddenPatterns = imported.ForbiddenPatterns
	host.MinCount = imported.MinCount
	host.MaxCount = imported.MaxCount
	host.MaxDepth = imported.MaxDepth
	host.AllowChildren = imported.AllowChildren
}

func isRemote(ref string) bool {
	return strings.HasPrefix(ref, "http://") || strings.HasPrefix(ref, "https://")
}

// canonicalize returns a stable key for cycle detection: the absolute
// filesystem path for local refs, or the URL as-is for remote refs.
// basePath is the file path of the importing template.
func canonicalize(ref string, basePath string) (string, error) {
	if isRemote(ref) {
		return ref, nil
	}
	resolved := ref
	if !filepath.IsAbs(ref) {
		resolved = filepath.Join(filepath.Dir(basePath), ref)
	}
	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", fmt.Errorf("resolve template path %q: %w", ref, err)
	}
	return abs, nil
}

// loadTemplate fetches a template (from disk or HTTP), parses it, and
// returns the parsed config plus the new basePath that further nested
// imports inside that template should be resolved against.
func loadTemplate(ref string, basePath string) (Config, string, error) {
	var data []byte
	var nextBase string

	if isRemote(ref) {
		resp, err := http.Get(ref)
		if err != nil {
			return Config{}, "", fmt.Errorf("fetch %q: %w", ref, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return Config{}, "", fmt.Errorf("fetch %q: status %d", ref, resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Config{}, "", fmt.Errorf("read body of %q: %w", ref, err)
		}
		data = body
		nextBase = ref // remote refs nested inside are resolved as absolute URLs
	} else {
		path := ref
		if !filepath.IsAbs(path) {
			path = filepath.Join(filepath.Dir(basePath), ref)
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return Config{}, "", fmt.Errorf("read template %q: %w", ref, err)
		}
		data = body
		nextBase = path
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, "", fmt.Errorf("decode imported template %q: %w", ref, err)
	}
	return cfg, nextBase, nil
}
