package mdcheck

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

type Options struct {
	Input           string
	JSON            bool
	NoDuplicates    bool
	AllowMultipleH1 bool
	MaxLevel        int
	GroupByRule     bool
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

	headings, err := markdown.ExtractHeadings(source)
	if err != nil {
		return 2, err
	}

	diags := Validate(headings, ValidateOptions{
		CheckDuplicates: !opts.NoDuplicates,
		AllowMultipleH1: opts.AllowMultipleH1,
		MaxLevel:        opts.MaxLevel,
	})

	if len(diags) == 0 {
		return 0, nil
	}

	output, err := renderDiagnostics(diags, opts.JSON, opts.GroupByRule)
	if err != nil {
		return 2, err
	}

	if output != "" {
		if _, err := fmt.Fprint(stdout, output); err != nil {
			return 2, err
		}
	}

	return 1, nil
}

func renderDiagnostics(diags []Diagnostic, asJSON bool, groupByRule bool) (string, error) {
	if asJSON {
		return RenderJSON(diags)
	}

	return RenderHuman(diags, groupByRule), nil
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	fs := flag.NewFlagSet("mdcheck", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.BoolVar(&opts.JSON, "json", false, "output diagnostics as JSON")
	fs.BoolVar(&opts.NoDuplicates, "no-duplicates", false, "disable duplicate heading check")
	fs.BoolVar(&opts.AllowMultipleH1, "allow-multiple-h1", false, "allow more than one H1 heading")
	fs.IntVar(&opts.MaxLevel, "max-level", 0, "maximum heading level allowed (0 disables the check)")
	fs.BoolVar(&opts.GroupByRule, "group-by-rule", false, "group diagnostics by rule, then by line number")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--json] [--group-by-rule] [--no-duplicates] [--allow-multiple-h1] [--max-level N] input.md\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if err := validateMaxLevel(opts.MaxLevel); err != nil {
		return opts, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input Markdown file")
	}

	opts.Input = fs.Arg(0)

	return opts, nil
}

func validateMaxLevel(maxLevel int) error {
	if maxLevel < 0 {
		return errors.New("max-level must be at least 0")
	}
	if maxLevel > 6 {
		return errors.New("max-level must be at most 6")
	}
	return nil
}
