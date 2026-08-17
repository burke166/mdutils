# CLAUDE.md

This file provides guidance for working in this repository.

## Project purpose

mdutils is a collection of small, focused command-line utilities for working
with Markdown files: outlining, generating tables of contents, splitting and
merging documents, formatting, linting, and computing document statistics.
Each tool does one job and is meant to be composed with others (shell
pipelines, CI checks, editor tooling).

## Language / module

Go, module path `github.com/computercodeblue/mdutils` (see `go.mod`),
targeting Go 1.26.

## Architecture

- `cmd/<name>/main.go` — one thin entrypoint per tool. Each just calls
  `Execute()` in the matching `internal/commands/<name>` package
  (e.g. `cmd/mdcheck/main.go` calls `mdcheck.Execute()`).
- `internal/commands/<name>/` — one package per tool, holding the tool's
  `Options`, `Execute`/`Run` functions, and any tool-specific logic
  (validation, rendering, splitting, merging, formatting, linting, stats).
- `internal/markdown/` — shared Markdown parsing: frontmatter extraction,
  heading extraction, line utilities. Built on `goldmark`.
- `internal/output/` — shared rendering helpers used across tools: JSON, CSV,
  tree, table of contents, bullet lists, numbered output, slug generation.
- `internal/fileutil/` — shared file walking (with skip-list for
  `.git`, `node_modules`, `vendor`, `bin`, `obj`, `dist`, `build`) and file
  writing helpers.
- `tools/build/` — a small Go program (not a CLI tool for end users) that
  globs `cmd/*`, builds each directory containing a `main.go`, and writes
  binaries to `bin/`. `--release` instead cross-compiles every tool for a
  fixed platform matrix (`windows/amd64`, `linux/amd64`, `linux/arm64`,
  `darwin/amd64`, `darwin/arm64`) into `dist/<goos>_<goarch>/`, with
  `CGO_ENABLED=0` — safe because the module has no cgo dependencies.
- `testdata/` — shared Markdown fixtures used by multiple tools' tests.
- `dist/` — placeholder for built distribution artifacts (kept in Git via
  `dist/.gitkeep`; contents are gitignored).

## Tools

| Tool | Purpose |
|---|---|
| `mdoutline` | Extract a document outline from a Markdown file |
| `mdtoc` | Generate a table of contents |
| `mdsplit` | Split a Markdown file into multiple files by heading level |
| `mdmerge` | Merge multiple Markdown files into one document |
| `mdformat` | Format Markdown files without changing their meaning |
| `mdlint` | Lint Markdown files for structure/style problems (read-only) |
| `mdstats` | Report document metrics (read-only) |
| `mdcheck` | Validate heading structure (single H1, no duplicates, etc.) |

See `README.md` for each tool's flags and example usage.

## Conventions already in use

- **CLI pattern**: every command package exposes
  `Run(args []string, stdout, stderr io.Writer) (int, error)` for testability,
  plus a thin `Execute()` that calls `Run(os.Args[1:], os.Stdout, os.Stderr)`
  and translates the result into `os.Exit`.
- Standard library `flag` package for argument parsing — no CLI framework
  dependency.
- No `context.Context` threading; these are short-lived, single-shot CLI
  invocations.
- Error handling is plain `error` returns; `fmt.Errorf` is used for
  wrapping/formatting where a caller needs more detail.
- Tests use `github.com/stretchr/testify/require`. Command-level tests
  reference shared fixtures in `testdata/` via relative paths
  (`filepath.Join("..", "..", "..", "testdata", ...)`).
- Read-only tools (`mdlint`, `mdstats`) document that they never modify
  files; destructive/mutating tools (`mdformat --write`, `mdsplit`,
  `mdmerge`) are explicit about when they write.
- Exit codes are meaningful and documented per tool in `README.md` (e.g.
  `mdlint`: 0 clean, 1 warnings, 2 errors, 3 internal failure) — preserve
  this convention when adding or changing exit behavior.

## Dependencies

Deliberately lightweight:

- `github.com/yuin/goldmark` — Markdown parsing.
- `github.com/stretchr/testify` — test assertions (dev/test only).
- `gopkg.in/yaml.v3` — frontmatter parsing.

Avoid adding new third-party dependencies unless a tool genuinely needs
functionality that isn't reasonably hand-rolled.

## Build / test / lint

```bash
# Build every cmd/* binary into bin/
go run ./tools/build
# or, via the wrapper scripts:
scripts/build-all.sh      # POSIX
scripts/build-all.ps1     # PowerShell

# Cross-compile every cmd/* binary for all release platforms into dist/
go run ./tools/build --release
scripts/build-all.sh --release
scripts/build-all.ps1 --release

# Compile-check everything
go build ./...

# Run the full test suite
go test ./...
```

There is no dedicated lint config (`.golangci.yml`) in the repo; run
`go vet ./...` as a baseline check before committing.

## Repository conventions

- `bin/` and `dist/` contents are gitignored; `dist/.gitkeep` is tracked so
  the directory survives in Git.
- Built binaries are Windows/Unix-portable via `tools/build`'s
  `executableName` helper (`.exe` suffix added on Windows).
- Keep `README.md`'s tool usage examples in sync with actual flag behavior
  when changing a command's CLI surface.

## Testing expectations

Every command package has a `command_test.go` exercising `Run` end-to-end
against fixtures, plus focused unit tests for its supporting logic (e.g.
`validate_test.go`, `split_test.go`, `stats_test.go`). New tools or flags
should follow the same shape: an end-to-end `Run` test plus unit tests for
any new non-trivial logic.

## Design-first workflow

Before starting non-trivial work (a new tool, a new flag with real behavior,
a change to shared `internal/` logic), write a short design doc in
`docs/design/` describing the problem and proposed approach. Keep it brief —
this is a small utility collection, not a large system — but capture the
reasoning that isn't obvious from the code itself.

## License

MIT (see `LICENSE`).
