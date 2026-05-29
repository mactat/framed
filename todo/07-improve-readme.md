# 07 — Improve README

**Status:** open
**Source:** added by maintainer (not in README's TODO list)

## Scope

**In scope**

- Reorganize the README so a new reader can get to "what is this and how do I try it" within 30 seconds.
- Trim duplication (especially the three near-identical install sections).
- Replace placeholders and broken references.
- Move the long-form TODO list out of the README and link to the [`todo/`](.) folder instead.

**Out of scope**

- Logo, badges, or other visual assets.
- A separate docs site (out of scope until/unless `link-to-full-docs` becomes real).
- Translating the README.

## Problems with the README today

Concrete issues, top to bottom of the current file:

1. **No table of contents** — file is ~320 lines and there's no nav.
2. **No Quick Start near the top** — usage examples don't appear until line ~200, after a long install section.
3. **Three duplicated install sections** (macOS / Linux / Windows) — they share the same 6 steps; only the filename and the final move-to-PATH step differ.
4. **Configuration reference is fragmented** — each YAML field gets its own H3 with one sentence (`### Forbidden Files`, `### Minimum File Count`, …). A single table would be denser and easier to scan.
5. **Inline 40-line example YAML** duplicates `examples/structure.yaml` — drift risk, and pushes the rest of the README down.
6. **Broken link** — `[documentation](link-to-full-docs)` is a placeholder.
7. **TODO list lives in the README** — now that `todo/` exists, the README should link to it instead of carrying the checklist inline.
8. **GitHub Action snippet pins `v0.0.8`** — easy to forget to bump; either parameterize or note that users should pin the latest release.

## Initial design

Proposed new section order:

1. **Header** — logo, badges, one-paragraph pitch.
2. **Quick Start** — three commands: install (brew), `framed import --example golang`, `framed verify`.
3. **Features** — bullet list (mostly as-is).
4. **Usage** — `create`, `verify`, `capture`, `visualize`, `import`. Keep examples short.
5. **Configuration** — one paragraph + a single reference table of every field (`files`, `dirs`, `allowedPatterns`, `forbiddenPatterns`, `minCount`, `maxCount`, `allowChildren`, `maxDepth`, `template:` once card [01](./01-nested-template-imports.md) lands). Link to `examples/structure.yaml` for the full example.
6. **Installation** — collapse the three OS sections into one "Manual install" block with a small table mapping OS → archive name → install path. Keep brew as its own short section.
7. **Docker** — as-is.
8. **GitHub Action** — as-is; replace the pinned version with a note to "use the latest release."
9. **Contributing** — short pointer to `CONTRIBUTING.md` (see card [03](./03-contributing-guidelines.md)) and `todo/`.

Field reference table sketch:

| Field | Applies to | Type | Description |
|---|---|---|---|
| `files` | dir | list | Required files in this directory |
| `dirs` | dir | list | Nested directories |
| `allowedPatterns` | dir | list | Glob patterns files must match |
| `forbiddenPatterns` | dir | list | Patterns files must not match |
| `minCount` / `maxCount` | dir | int | File count bounds |
| `allowChildren` | dir | bool | Allow nested subdirectories (default true) |
| `maxDepth` | dir | int | Max nesting depth |

### Acceptance criteria

- New reader can run their first `framed verify` from a fresh repo after reading the top of the README only.
- No duplicated install steps across OSes.
- No placeholder links.
- README is meaningfully shorter than today (target: under 220 lines).
- TODO checklist is gone from the README; replaced by a one-line pointer to `todo/`.

### Dependencies

- Card [03](./03-contributing-guidelines.md) — if `CONTRIBUTING.md` lands first, the new "Contributing" pointer in the README has somewhere to link.
- Card [01](./01-nested-template-imports.md) — when `template:` ships, add a row to the field reference table.
