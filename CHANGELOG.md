# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `.github/workflows/release.yml` — pushing a `v*.*.*` tag now automatically
  runs the test suite, cross-compiles release binaries, packages them per
  platform, and publishes a GitHub Release with the archives attached.

## [0.1.1] - 2026-08-17

### Added

- `tools/build --release` — cross-compile every tool for `windows/amd64`,
  `linux/amd64`, `linux/arm64` (Raspberry Pi, AWS Graviton), `darwin/amd64`,
  and `darwin/arm64` into `dist/<goos>_<goarch>/`. The default
  `go run ./tools/build` still builds only for the host platform into `bin/`.

### Documentation

- Document `mdcheck` in `README.md` (it was implemented and tested but
  missing from the tool list).

## [0.1.0] - 2026-08-17

### Added

- `mdoutline` — extract a document outline from a Markdown file, with
  `--tree` and `--json` output modes.
- `mdtoc` — generate a table of contents for a Markdown file, with
  `--min-level`/`--max-level` heading range control.
- `mdsplit` — split a Markdown file into multiple files by heading level,
  with `--level`, `--out`, and `--numbered` (zero-padded sequential file
  naming) options.
- `mdmerge` — merge multiple Markdown files into one document in the order
  given (or directory order, aided by numeric filename prefixes).
- `mdformat` — format Markdown files without changing their meaning, with
  `--write` and `--check` modes.
- `mdlint` — lint Markdown files for structure and style problems
  (single-H1, no-skipped-heading-levels, no-duplicate-headings,
  no-empty-headings/sections, whitespace, line length, and more), configurable
  via `.mdlintrc.yaml`, with `--json`, `--quiet`, and `--no-recursive` options
  and meaningful exit codes.
- `mdstats` — analyze Markdown files and report document metrics, with
  `--json`, `--csv`, `--summary`, `--exclude`, and `--no-recursive` options.
- `mdcheck` — validate heading structure (duplicate headings, multiple H1s,
  maximum heading level), with `--json` and `--group-by-rule` output options.
- `tools/build` and `scripts/build-all.{sh,ps1}` — build all `cmd/*` binaries
  into `bin/` in one step.

[Unreleased]: https://github.com/burke166/mdutils/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/burke166/mdutils/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/burke166/mdutils/releases/tag/v0.1.0
