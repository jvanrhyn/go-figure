# GitHub Copilot Instructions — go-figure

## Branch Workflow

Always do user-requested work on a feature branch, never directly on `master`.

- If the current branch is `master` (or another default/protected branch), create and switch to a feature branch before making changes, running generators, committing, or opening a PR.
- If the user names a branch, use that exact name.
- If the user does not name a branch, create a descriptive branch name from the task, using a prefix like `feature/`, `fix/`, `chore/`, `docs/`, `refactor/`, or `test/`.
- If already on a non-default branch, stay on it unless the user asks to change branches.
- Do not merge directly into `master` as part of fulfilling a request unless the user explicitly asks for that.
- Before any potentially branch-sensitive operation, verify the current branch with `git branch --show-current`.

## Commands

```bash
go test -race ./...              # full test suite
go test -race -run TestNewFigure # single test
go vet ./...                     # vet
golangci-lint run                # lint (requires golangci-lint installed)
gofmt -w .                       # format (always run before committing)
```

## Architecture

Single-package library (`package figure`, module `github.com/common-nighthawk/go-figure`) with no external dependencies. Five source files:

| File | Responsibility |
|------|---------------|
| `figure.go` | Exported `Figure` struct, three constructors, `Slicify`, `reverse` helper |
| `public_methods.go` | Methods on `Figure`: `Print`, `String`, `ColorString`, `Write`, `Scroll`, `Blink`, `Dance` |
| `font.go` | Unexported `font` struct; loads `.flf` files from embedded FS (`//go:embed fonts`); ANSI color map |
| `figlet-parser.go` | Parses FIGlet file header fields; all functions take a pre-split `[]string` and return errors |
| `figure_test.go` | 18 tests covering constructors, parser helpers, error paths, and `Write` |

149 bundled `.flf` font files live in `fonts/` and are embedded at compile time via `embed.FS` in `font.go`. No generated code.

## Key Conventions

**All constructors and `Slicify` return errors** — this is a breaking change from the original upstream API:
```go
fig, err := figure.NewFigure("Hello", "standard", true)
rows, err := fig.Slicify()
```

**`font` is unexported and embedded in `Figure`** — don't export it. The public surface is `Figure` only.

**`strict` mode**: when `true`, non-ASCII runes return an error; when `false`, they are silently replaced with `?`.

**FIGlet parser inputs are pre-split**: `figlet-parser.go` functions accept `[]string` (already `strings.Fields`-split header) and perform bounds checks before indexing. Any new parser helpers must follow this pattern and return `(value, error)`.

**`gofmt` is intentionally absent from `.golangci.yml`** — golangci-lint is built with an older Go toolchain; adding `gofmt` causes version-skew failures. Always run `gofmt -w .` manually before committing; correctness linters (`govet`, `staticcheck`, `errcheck`, `unused`) still run in CI.

**Example test functions** use the `Example_suffix` form (not `ExampleSymbol`) because the examples demonstrate package-level usage, not a single exported symbol:
```go
func Example_standard() { ... }
func Example_alphabet() { ... }
```

**Error wrapping in `font.go`**: use `fmt.Errorf("... : %w", err)` to wrap underlying parse errors so callers can use `errors.Is`/`errors.As`.

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on every push and PR:
- Matrix: `{stable, oldstable}` × `{ubuntu-latest, macos-latest}`
- Steps: `go vet` → golangci-lint → `go test -race ./...`

The `go` directive in `go.mod` is pinned to `1.21` (minimum for the hard-minimum semantics); the actual minimum runtime requirement is Go 1.16 (`embed` package).
