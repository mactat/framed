# 03 — Contributing guidelines

**Status:** open
**Source:** README "TODO"

## Scope

**In scope**

- A `CONTRIBUTING.md` at repo root covering the basics a new contributor needs.
- A PR template under `.github/PULL_REQUEST_TEMPLATE.md`.
- An issue template (bug + feature) under `.github/ISSUE_TEMPLATE/`.
- A short "Development" section linked from the README.

**Out of scope**

- Code of Conduct (separate file/effort).
- Governance / maintainer roster.
- Release process documentation (already partly captured by the release workflow).

## Initial design

`CONTRIBUTING.md` outline:

1. **Project layout** — short tour of `cmd/`, `pkg/ext/`, `examples/`, `test/`.
2. **Local development** — Go version (from `go.mod`), `make build`, `make test`, running the binary against `framed.yaml`.
3. **Tests** — how to run `test.bats` locally (bats prerequisite), how to add a new fixture.
4. **Style** — `gofmt` / `go vet`; commit message style (look at recent log: short imperative subject).
5. **PR checklist** — tests pass, README/docs updated if behavior changes, `framed verify` passes.
6. **Reporting bugs** — link to issue template.

PR template fields: summary, motivation, testing notes, related issue.

Issue templates:

- `bug_report.md` — repro steps, expected vs actual, version (`framed --version` if available).
- `feature_request.md` — problem, proposed solution, alternatives considered.

### Dependencies

None — pure docs. Can ship independently of other open cards.
