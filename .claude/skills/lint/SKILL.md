---
name: lint
description: Run gofmt, go vet, and golangci-lint against this repo, then surface a prioritized list of findings. Use when the user asks to "lint", "check formatting", "run golangci-lint", or chases CI lint failures.
---

# Lint

Runs the same checks CI runs (`gofmt`, `go vet`, `golangci-lint`) and
reports findings. Used to chase CI lint failures locally without
waiting for a CI round-trip.

## When to use this skill

- The user asks to "lint", "run lint", "check formatting".
- CI is failing on the `golangci-lint` job and the user wants a
  local reproduction.
- Before opening a PR, when the user wants to know what CI will say.

## How to use this skill

1. **Confirm the binary is available.** `golangci-lint --version`.
   If it's missing, install the version pinned in
   `.github/workflows/pr.yaml` (currently `v1.53`):

   ```shell
   curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
     | sh -s -- -b /usr/local/bin v1.53.3
   ```

   On modern Go versions `go install` of v1.53.3 fails to compile —
   use the install script above.

2. **Warm the module cache.** golangci-lint's typecheck pass reports
   `undefined: yaml` / `undefined: color` when modules aren't
   resolved. Run first:

   ```shell
   go mod download
   ```

3. **Run the same checks CI runs**, in this order:

   ```shell
   gofmt -l .                       # lists unformatted files
   go vet ./...
   golangci-lint run -v
   ```

   If running through Docker (matches CI exactly):

   ```shell
   make lint
   ```

4. **Triage the findings.** Group by file, then prioritize:
   - Compile errors (`typecheck`) first — they mask everything else.
   - Bugs (`govet`, `errcheck`, `ineffassign`) next.
   - Style (`gofmt`, `unused`, `gosimple`) last.

5. **Fix and re-run.** Don't claim done until `golangci-lint run`
   exits 0. `go fmt ./...` (or `make format`) handles formatting in
   bulk.

## Common gotchas

- **"undefined: yaml" / "undefined: color"**: dependencies not
  downloaded. Always `go mod download` before linting.
- **Lint version drift**: CI pins `v1.53` in `.github/workflows/pr.yaml`.
  Newer golangci-lint versions surface different findings — use the
  pinned version locally.
- **CGO_ENABLED**: this project builds fine with `CGO_ENABLED=0`. If
  your local gcc is missing flags, set it before invoking go tools.
