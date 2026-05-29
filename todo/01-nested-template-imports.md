# 01 — Nested template imports

**Status:** open
**Source:** README "TODO" — first item

## Scope

**In scope**

- Extend the YAML schema so any `dirs:` entry may declare a `template:` field that points to either:
  - a local path (`./other.yaml`, `../shared/lib.yaml`), or
  - a remote URL (`https://...framed.yaml`).
- When a `template:` is present on a dir, its contents replace (or are merged into) the dir definition during config load.
- Support nesting: an imported template may itself reference further templates.
- Surface clear errors for: missing files, non-2xx HTTP responses, malformed YAML, cycles.

**Out of scope (for first cut)**

- Authenticated URLs (tokens, private repos).
- Caching of remote templates.
- Schema versioning / migration.
- Partial-template merging semantics beyond "replace the dir subtree."

## Initial design

Loader changes live in `pkg/ext/decoder.go`. The existing `ReadConfig` returns a tree plus a flat dir list; we add a resolution pass between parse and flatten.

```
ReadConfig(path)
  └─ parseYAML(path) → tree
  └─ resolveImports(tree, baseDir, seen)   ← NEW
       for each node with `template:`
         resolved := loadTemplate(ref, baseDir, seen)   // local file or URL
         node.replaceWith(resolved.structure)
  └─ flatten(tree) → dirList
```

- `loadTemplate` reuses `ext.ImportFromUrl` for the URL case (already in `pkg/ext/network.go`) and `os.ReadFile` for the local case. URL detection: `strings.HasPrefix(ref, "http://" | "https://")`.
- `seen` is a `map[string]bool` keyed by absolute path / canonical URL, populated before recursion to detect cycles.
- Relative local paths are resolved against the *importing* template's directory, not the CWD.
- Merge rule for v1: the imported template's `structure` block fully replaces the host dir's children. Host dir keeps its `name`. Field-level merging can come later if requested.

### YAML example

```yaml
name: framed
structure:
  name: root
  dirs:
    - name: shared
      template: ./shared.yaml
    - name: remote
      template: https://example.com/framed.yaml
```

### Testing

Add bats cases under `test/` covering: local import, remote import (httptest server or fixture URL), cycle detection, missing file, bad URL. Reuse the existing `test.bats` harness.

### Open questions

- Should `template:` and inline `dirs:`/`files:` on the same node be allowed (merge) or rejected? Start with **rejected** to keep semantics simple.
- For URLs, should we pin by hash for reproducibility? Out of scope for v1, worth a follow-up.
