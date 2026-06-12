package mdtoc

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/computercodeblue/mdutils/internal/output"
)

type Options struct {
	Input    string
	MinLevel int
	MaxLevel int
	Ordered  bool
	NoLinks  bool
	NoIndent bool
}

func Execute() {
	if err := Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	opts, err := parseOptions(args, stderr)
	if err != nil {
		return err
	}

	source, err := os.ReadFile(opts.Input)
	if err != nil {
		return err
	}

	headings, err := markdown.ExtractHeadings(source)
	if err != nil {
		return err
	}

	toc := output.RenderToc(headings, output.TocOptions{
		MinLevel: opts.MinLevel,
		MaxLevel: opts.MaxLevel,
		Ordered:  opts.Ordered,
		NoLinks:  opts.NoLinks,
		NoIndent: opts.NoIndent,
	})

	_, err = fmt.Fprint(stdout, toc)
	return err
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	opts.MinLevel = 1
	opts.MaxLevel = 6

	fs := flag.NewFlagSet("mdtoc", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.IntVar(&opts.MinLevel, "min-level", 1, "minimum heading level to include")
	fs.IntVar(&opts.MaxLevel, "max-level", 6, "maximum heading level to include")
	fs.BoolVar(&opts.Ordered, "ordered", false, "use ordered list markers")
	fs.BoolVar(&opts.NoLinks, "no-links", false, "output plain heading text without links")
	fs.BoolVar(&opts.NoIndent, "no-indent", false, "output without indenting nested items")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--min-level N] [--max-level N] [--ordered] [--no-links] [--no-indent] input.md\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if err := validateLevels(opts.MinLevel, opts.MaxLevel); err != nil {
		return opts, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input Markdown file")
	}

	opts.Input = fs.Arg(0)

	return opts, nil
}

func validateLevels(min, max int) error {
	if min < 1 {
		return errors.New("min-level must be at least 1")
	}
	if max > 6 {
		return errors.New("max-level must be at most 6")
	}
	if min > max {
		return errors.New("min-level must not be greater than max-level")
	}
	return nil
}
