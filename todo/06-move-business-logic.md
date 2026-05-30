# 06 — Move remaining business logic to a separate package

**Status:** done
**Source:** README "TODO" — last item

## Scope

**In scope**

- Lift verification / capture / create / visualize logic out of `cmd/*.go` into `pkg/ext/` (or a new sub-package).
- `cmd/*.go` becomes a thin Cobra layer: parse flags → call into the package → render result / set exit code.
- No behavior change; pure refactor.

**Out of scope**

- Renaming the binary or CLI surface.
- Changing YAML schema or output format.
- Splitting `pkg/ext` into finer-grained packages (can be a follow-up if `ext` gets too broad).

## Initial design

### Today

- `cmd/verify.go` (93 lines) inlines the verification loop — it walks the dir list, calls helpers in `pkg/ext`, and prints/exits directly.
- `cmd/capture.go`, `cmd/create.go`, `cmd/visualize.go`, `cmd/import.go` follow the same shape: flag parsing mixed with orchestration logic and `os.Exit` calls.
- Reusable primitives (`DirExists`, `CountFiles`, `VerifyFiles`, `CheckDepth`, ...) already live in `pkg/ext/system.go` — good. The *orchestration* on top of them is what's still in `cmd/`.

### Target shape

For each command, extract a single entry-point function into `pkg/ext` that returns a structured result instead of exiting:

```go
// pkg/ext/verify.go
type VerifyReport struct {
    Failures []string
    OK       bool
}
func Verify(templatePath string) (VerifyReport, error) { ... }
```

```go
// cmd/verify.go
Run: func(cmd *cobra.Command, args []string) {
    report, err := ext.Verify(cmd.Flag("template").Value.String())
    if err != nil { ext.PrintOut("☠️", err); os.Exit(2) }
    for _, f := range report.Failures { ext.PrintOut("❌", f) }
    if !report.OK { os.Exit(1) }
    fmt.Println("✅ Verified successfully!")
},
```

Same treatment for `Capture`, `Create`, `Visualize`, `Import`.

### Why returns instead of prints

- Makes the logic unit-testable (today it's only reachable end-to-end via bats).
- Keeps `os.Exit` confined to `cmd/`, which is the only layer that should own process lifecycle.
- Sets up the codebase to be embedded as a library (someone could `import "framed/pkg/ext"` and call `Verify` from their own tool).

### Migration order

1. `Verify` first — smallest blast radius, most-tested behavior, biggest unit-test win.
2. `Create` and `Capture` next.
3. `Visualize` and `Import` last.

Each step is a separate PR; behavior pinned by the existing bats suite ([card 02](./02-add-tests.md)).

### Dependencies

Unblocks easier unit testing (extends [card 02](./02-add-tests.md)). No hard dependency on other cards.
