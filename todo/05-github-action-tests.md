# 05 — GitHub Action for running tests

**Status:** done
**Source:** README "TODO" — marked `[x]`

## Scope

CI already runs the test suite. Workflows live in `.github/workflows/`:

- `pr.yaml` — runs on PRs (build + tests).
- `release.yaml` — runs on tags (build + publish).

The README's PR/release badges at the top reference these.

## Initial design

Already shipped. Follow-up ideas if/when wanted:

- Add a `lint` job (`gofmt -l`, `go vet`).
- Cache the Go module download to speed up runs.
- Run `framed verify` against the repo itself in CI as a dogfooding check (the repo has a `framed.yaml`).
- Matrix-build on multiple Go versions if compatibility becomes a concern.
