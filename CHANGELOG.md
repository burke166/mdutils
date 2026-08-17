# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/burke166/mdutils/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/burke166/mdutils/releases/tag/v0.1.0
