# 04 — More examples

**Status:** done
**Source:** README "TODO"

## Scope

**In scope**

- New example templates under `examples/`, each a realistic project layout that uses several Framed features (`allowedPatterns`, `minCount`/`maxCount`, `allowChildren`, `maxDepth`, nested `dirs:`).
- Each example wired into the `import --example <name>` flow (see `ext.ExampleToUrl` in `pkg/ext/network.go`).
- README "Currently available examples" list updated.

**Candidate additions** (pick a subset for the first batch):

- `node` — typical Node/TypeScript service (`src/`, `tests/`, `package.json`, `tsconfig.json`).
- `rust` — Cargo layout (`src/`, `tests/`, `benches/`, `Cargo.toml`).
- `terraform` — module layout (`modules/`, `environments/`, `*.tf`).
- `monorepo` — multi-service layout exercising `maxDepth` and nested dirs.
- `docs-site` — `docs/`, `static/`, MkDocs/Hugo-style.

**Out of scope**

- Language-specific scaffolding beyond the directory shape (no source code in examples — Framed only cares about structure).

## Initial design

Each new example is one YAML file in `examples/`, mirroring the shape of the existing `golang.yaml` / `python.yaml`. Steps per example:

1. Author `examples/<name>.yaml`.
2. Confirm `ExampleToUrl("<name>")` resolves correctly (it builds a URL from the example name; check the format in `pkg/ext/network.go`).
3. Add a bats case: `framed import --example <name> --output /tmp/x.yaml` succeeds and the output parses.
4. Update README's "Currently available examples" list.

### Acceptance criteria per example

- Loads via `framed import --example <name>`.
- Passes `framed verify` against a freshly `framed create`d directory built from it.
- Demonstrates at least one non-trivial rule (pattern, count, depth) — otherwise it's just a file list.

### Dependencies

Loosely related to [card 01](./01-nested-template-imports.md): once nested imports land, a "monorepo" example could compose smaller per-service templates instead of duplicating them.
