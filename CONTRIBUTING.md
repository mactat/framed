# Contributing to Framed

Thanks for considering a contribution. Framed is a small Go CLI; the
contributing surface is correspondingly small. This document covers
what you need to know to land a change.

## Project layout

```
cmd/        Cobra command wrappers. Thin: flag parsing → call into pkg/ext → exit.
pkg/ext/    All business logic — decoder, encoder, system probes, verify, capture,
            create, visualize, import. Unit tests live next to the code (*_test.go).
examples/   Sample YAML templates (golang, python, node, rust, terraform, monorepo,
            docs-site). Every example must round-trip through create + verify.
test/       bats-based end-to-end CLI tests. test_helper/ is fetched at build time
            and gitignored.
dockerfiles/ Build, release, and test container definitions.
todo/       One markdown card per planned work item.
```

## Local development

1. Install Go (see `go.mod` for the required version).
2. Build the binary:

   ```shell
   make build
   ```

   The output goes to `./build/framed`. Put it on your `PATH` if you
   want to use it locally:

   ```shell
   export PATH=$PWD/build:$PATH
   ```

3. Drive the binary against the repo's own template:

   ```shell
   framed verify --template ./framed.yaml
   ```

## Running tests

Two test layers, both required:

**Go unit tests** (fast, pure-Go):

```shell
go test ./...
```

**Bats end-to-end tests** (driver: the built binary):

Prerequisites: `bats`, `yq`, and the bats helper libraries cloned under
`test/test_helper/`. The Docker-based path handles this for you:

```shell
make test
```

To run bats directly without Docker:

```shell
# one-time setup
mkdir -p test/test_helper
git clone https://github.com/bats-core/bats-support.git  test/test_helper/bats-support
git clone https://github.com/bats-core/bats-assert.git   test/test_helper/bats-assert
git clone https://github.com/bats-core/bats-file.git     test/test_helper/bats-file

# build, install on PATH, run
make build && sudo cp build/framed /usr/local/bin/
bats test/*.bats
```

## Adding a new test fixture

For bats: drop a `<name>.bats` file in `test/`. Use the existing
`test/test.bats` as a template — call `load
'test_helper/bats-support/load'` etc. in `setup()`, isolate to a
`/tmp/...` working dir, and clean up in `teardown()`.

For Go: drop a `<thing>_test.go` next to the code it tests in
`pkg/ext/`. Use `t.TempDir()` for filesystem fixtures and the
`runIn()` / `writeTemplate()` helpers in `verify_test.go` for
templates that need a particular CWD.

## Style

- `gofmt` everything (`make format`).
- `go vet ./...` should pass.
- `golangci-lint` is run in CI (`make lint`).
- Commit messages: short, imperative subject (look at `git log` for
  the prevailing style).

## Pull request checklist

Before opening a PR:

- [ ] `go test ./...` passes.
- [ ] `bats test/*.bats` passes (or `make test`).
- [ ] If behavior changes, the README is updated.
- [ ] If a new feature is added, an example or a `todo/` card covers
      it.
- [ ] `framed verify` still passes on the repo itself.

## Reporting bugs and requesting features

Use the issue templates under `.github/ISSUE_TEMPLATE/`. They prompt
for the information that's almost always needed (repro steps,
expected vs actual, version).

## Planning

Open work is tracked as cards under [`todo/`](./todo/). Each card has a
status, a scope, and an initial design sketch. If you're about to
start a non-trivial change, look there first — and if you're proposing
something new, consider adding a card before writing code so the
design can be discussed.
