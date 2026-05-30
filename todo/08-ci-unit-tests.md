# 08 — Extend CI to run Go unit tests

**Status:** done
**Source:** follow-up to [card 06](./06-move-business-logic.md) and
[card 02](./02-add-tests.md)

## Scope

**In scope**

- Make `go test ./...` a required check on every PR.
- Run it in `.github/workflows/pr.yaml` as a separate job (or a
  step) from the existing bats job, so a Go-level regression is
  visible without waiting for the Docker-built bats run.
- Cache `~/go/pkg/mod` so consecutive runs don't repay the download
  cost.
- Surface coverage somewhere a human will actually see it — at
  minimum, print the per-package `-cover` summary in the job log.

**Out of scope**

- Coverage gating (requiring N% coverage before merge).
- Matrix-build across multiple Go versions.
- Race-detector runs (`go test -race`) — track separately if needed.
- Coverage upload to a third-party service (Codecov, Coveralls).

## Why

Today CI only runs the bats end-to-end suite. The Go unit tests
under `pkg/ext/*_test.go` (added in card 06 and extended in cards 01,
04) are not exercised by CI at all — a developer can land a PR that
breaks them and only notice when they next run `go test` locally.

Symptoms:

- Quick feedback loop is missing. Bats requires the Docker image to
  be built before any test runs (~minutes); `go test ./...` is
  seconds.
- A failure in a unit test won't block merge even though it would
  block local development.

## Initial design

### Today

`.github/workflows/pr.yaml` runs:

1. `mactat/framed-action@v0.0.7` — dogfooded `framed verify`.
2. `make test BUILD=true EXPORT=true` — bats via Docker.
3. `test-summary/action@v2` — bats result rendering.
4. `golangci-lint-action@v3` — already runs.

There is no `go test` step.

### Target shape

Add a `unit` job that runs alongside the existing `pr` job:

```yaml
unit:
  name: "Go unit tests"
  runs-on: "ubuntu-latest"
  steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version-file: go.mod
        cache: true            # caches ~/go/pkg/mod automatically
    - run: go test ./... -cover
```

Two-job split (instead of one combined job) so the unit tests
post their own status check and a bats regression doesn't mask a
unit-test failure.

### Dependencies

- Card [06](./06-move-business-logic.md) — needs unit tests to exist,
  which they now do.
- Card [05](./05-github-action-tests.md) — already shipped the
  workflow file that this card extends.

### Open questions

- Worth caching by `go.sum` hash explicitly, or is
  `setup-go cache: true` enough? Start with the latter; revisit if
  cache misses become noticeable.
- Add `-race`? Worth it for the import HTTP-server tests, but adds
  ~30s per run. Defer to a follow-up unless flakes appear.
