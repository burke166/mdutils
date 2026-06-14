# mdutils

A collection of small, focused command-line utilities for working with Markdown files.

## Building

```bash
go run ./tools/build
```

Binaries are written to `bin/`.

## Tools

### mdoutline

Extract a document outline from a Markdown file. Useful for documents converted to Markdown from PDF or DOCX that will require further editing.

```bash
mdoutline README.md
mdoutline --tree README.md
mdoutline --json README.md
```

### mdtoc

Generate a table of contents for a Markdown file.

```bash
mdtoc README.md
mdtoc --min-level 2 --max-level 4 README.md
```

### mdsplit

Split a Markdown file into multiple files by heading level.

```bash
mdsplit handbook.md
mdsplit handbook.md --numbered
mdsplit handbook.md --level 2 --numbered
mdsplit --level 2 --out parts README.md
```

With `--numbered` (or `-n`), each output file is prefixed with a zero-padded sequential number so files sort correctly in file managers and Git:

```text
01_introduction.md
02_installation.md
03_configuration.md
04_examples.md
05_reference.md
```

### mdmerge

Merge multiple Markdown files into one document. Files will be added in the order they're listed in the directory. Appending a prefix like 01_, 02_, etc. will ensure your files will be assembled in the order you want.

```bash
mdmerge chapter-*.md
mdmerge docs/ --out combined.md
```

### mdformat

Format Markdown files without changing their meaning.

```bash
mdformat README.md
mdformat --write README.md
mdformat --check README.md
```

### mdlint

Lint Markdown files for common structure and style problems. Reports issues but does not modify files.

```bash
# Lint a single file
mdlint README.md

# Lint the current directory recursively
mdlint .

# Lint a folder without recursing into subdirectories
mdlint --no-recursive docs/

# Use a specific config file
mdlint --config .mdlintrc.yaml .

# Machine-readable JSON output
mdlint --json README.md

# Only print errors (warnings are still counted in the exit code)
mdlint --quiet .
```

#### Configuration

`mdlint` reads `.mdlintrc.yaml` from:

1. The current working directory
2. The user's home directory
3. The executable directory

Use `--config <path>` to load a specific config file instead.

Example `.mdlintrc.yaml`:

```yaml
rules:
  single-h1: true
  no-missing-h1: true
  no-skipped-heading-levels: true
  no-duplicate-headings: true
  no-empty-headings: true
  no-empty-sections: true
  no-trailing-whitespace: true
  max-heading-length: 80
  max-line-length: 120
  no-multiple-blank-lines: true
  require-code-fence-language: false
  no-empty-links: true
  require-image-alt-text: false

severity:
  single-h1: error
  no-missing-h1: warning
  no-skipped-heading-levels: warning
  no-duplicate-headings: warning
  no-empty-headings: error
  no-empty-sections: warning
  no-trailing-whitespace: warning
  max-heading-length: warning
  max-line-length: warning
  no-multiple-blank-lines: warning
  require-code-fence-language: warning
  no-empty-links: error
  require-image-alt-text: warning

exclude:
  - "CHANGELOG.md"
  - "docs/generated/**"
```

#### Exit codes

| Code | Meaning |
|------|---------|
| 0 | No issues found |
| 1 | Warnings only |
| 2 | One or more errors |
| 3 | Internal failure |

#### Example output

```text
README.md
  line 1: error single-h1: document has multiple H1 headings
  line 18: warning no-skipped-heading-levels: heading level skipped from H2 to H4
  line 55: warning no-duplicate-headings: duplicate heading "Installation"

3 issues found
```

When linting folders, `mdlint` skips common generated/vendor directories (`.git`, `node_modules`, `vendor`, `bin`, `obj`, `dist`, `build`).

### mdstats

Analyze Markdown files and report document metrics. Read-only; does not modify files.

```bash
# Analyze a single file
mdstats README.md

# Analyze the current directory recursively
mdstats .

# Analyze a folder without recursing into subdirectories
mdstats --no-recursive docs/

# Machine-readable JSON output
mdstats --json README.md

# Tabular CSV output (one row per file)
mdstats --csv .

# Combined summary for all files
mdstats --summary --per-file=false .

# Exclude files by glob pattern
mdstats --exclude CHANGELOG.md --exclude "docs/generated/**" .
```

#### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Analysis succeeded |
| 2 | One or more files could not be read |
| 3 | Internal failure |

When analyzing folders, `mdstats` skips common generated/vendor directories (`.git`, `node_modules`, `vendor`, `bin`, `obj`, `dist`, `build`).

## License

MIT. See [LICENSE](LICENSE).
