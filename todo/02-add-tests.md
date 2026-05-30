# 02 — Add some tests

**Status:** done
**Source:** README "TODO" — marked `[x]`

## Scope

A baseline test suite exists. Current state:

- `test/test.bats` — bats-based end-to-end CLI tests.
- `test/test.yaml` — fixture template used by those tests.
- CI runs them on PR (see [card 05](./05-github-action-tests.md)).

## Initial design

Already shipped. Follow-up work tracked separately if needed:

- Add Go-level unit tests under `pkg/ext/*_test.go` (decoder, encoder, system, network). Currently zero `_test.go` files exist.
- Expand bats coverage for error paths and exit codes (`verify` failure modes, malformed templates, missing files).
- Consider table-driven tests once business logic moves out of `cmd/` (see [card 06](./06-move-business-logic.md)) — that refactor makes the logic easier to unit-test.
