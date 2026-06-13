package mdformat

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/computercodeblue/mdutils/internal/fileutil"
)

type Options struct {
	Input               string
	Output              string
	Write               bool
	Check               bool
	TrimTrailingSpace   bool
	EnsureFinalNewline  bool
	NormalizeHeadings   bool
	NormalizeLists      bool
	PreserveFrontmatter bool
}

func Execute() {
	code, err := Run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	os.Exit(code)
}

func Run(args []string, stdout io.Writer, stderr io.Writer) (int, error) {
	opts, err := parseOptions(args, stderr)
	if err != nil {
		return 2, err
	}

	source, err := os.ReadFile(opts.Input)
	if err != nil {
		return 2, err
	}

	formatted := FormatMarkdown(string(source), FormatOptions{
		TrimTrailingSpace:   opts.TrimTrailingSpace,
		EnsureFinalNewline:  opts.EnsureFinalNewline,
		NormalizeHeadings:   opts.NormalizeHeadings,
		NormalizeLists:      opts.NormalizeLists,
		PreserveFrontmatter: opts.PreserveFrontmatter,
	})

	if opts.Check {
		if bytes.Equal(source, []byte(formatted)) {
			return 0, nil
		}
		return 1, nil
	}

	formattedBytes := []byte(formatted)

	switch {
	case opts.Write:
		if bytes.Equal(source, formattedBytes) {
			return 0, nil
		}
		if err := fileutil.WriteFileAtomically(opts.Input, formattedBytes); err != nil {
			return 2, err
		}
	case opts.Output != "":
		if err := fileutil.WriteFileAtomically(opts.Output, formattedBytes); err != nil {
			return 2, err
		}
	default:
		if _, err := fmt.Fprint(stdout, formatted); err != nil {
			return 2, err
		}
	}

	return 0, nil
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	opts.TrimTrailingSpace = true
	opts.EnsureFinalNewline = true
	opts.NormalizeHeadings = true
	opts.NormalizeLists = true
	opts.PreserveFrontmatter = true

	fs := flag.NewFlagSet("mdformat", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&opts.Write, "write", false, "rewrite the file in place")
	fs.BoolVar(&opts.Check, "check", false, "exit 0 if already formatted, 1 if changes would be made")
	fs.StringVar(&opts.Output, "output", "", "write formatted output to a different file")
	fs.BoolVar(&opts.TrimTrailingSpace, "trim-trailing-space", true, "remove trailing whitespace")
	fs.BoolVar(&opts.EnsureFinalNewline, "ensure-final-newline", true, "ensure the file ends with exactly one newline")
	fs.BoolVar(&opts.NormalizeHeadings, "normalize-headings", true, "normalize ATX heading spacing")
	fs.BoolVar(&opts.NormalizeLists, "normalize-lists", true, "normalize unordered list markers to \"-\"")
	fs.BoolVar(&opts.PreserveFrontmatter, "preserve-frontmatter", true, "preserve YAML frontmatter exactly as written")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--write] [--check] [--output file] [options] file.md\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if opts.Write && opts.Output != "" {
		return opts, errors.New("choose only one of --write or --output")
	}
	if opts.Check && opts.Write {
		return opts, errors.New("choose only one of --check or --write")
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input Markdown file")
	}

	opts.Input = fs.Arg(0)

	if opts.Output != "" {
		inputPath, err := filepath.Abs(opts.Input)
		if err != nil {
			return opts, err
		}
		outputPath, err := filepath.Abs(opts.Output)
		if err != nil {
			return opts, err
		}
		if inputPath == outputPath {
			return opts, errors.New("--output must differ from the input file")
		}
	}

	return opts, nil
}
