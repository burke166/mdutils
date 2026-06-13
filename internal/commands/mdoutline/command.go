package mdoutline

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/computercodeblue/mdutils/internal/fileutil"
	"github.com/computercodeblue/mdutils/internal/markdown"
)

type Options struct {
	Format string
	Input  string
	Output string
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

	output, err := Render(headings, opts.Format)
	if err != nil {
		return 2, err
	}

	if opts.Output != "" {
		if err := fileutil.WriteFileAtomically(opts.Output, []byte(output)); err != nil {
			return 2, err
		}
		return 0, nil
	}

	if _, err := fmt.Fprint(stdout, output); err != nil {
		return 2, err
	}

	return 0, nil
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	fs := flag.NewFlagSet("mdoutline", flag.ContinueOnError)
	fs.SetOutput(stderr)

	bullets := fs.Bool("bullets", false, "output bullet outline")
	tree := fs.Bool("tree", false, "output tree outline")
	numbered := fs.Bool("numbered", false, "output numbered outline")
	jsonOut := fs.Bool("json", false, "output JSON")
	csvOut := fs.Bool("csv", false, "output CSV")
	headingsOut := fs.Bool("headings", false, "output Markdown headings")

	fs.StringVar(&opts.Output, "output", "", "output file")
	fs.StringVar(&opts.Output, "o", "", "output file")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--bullets|--tree|--numbered|--json|--csv|--headings] [-o file] input.md\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	formats := map[string]bool{
		"bullets":  *bullets,
		"tree":     *tree,
		"numbered": *numbered,
		"json":     *jsonOut,
		"csv":      *csvOut,
		"headings": *headingsOut,
	}

	for name, used := range formats {
		if !used {
			continue
		}

		if opts.Format != "" {
			return opts, errors.New("choose only one output format")
		}

		opts.Format = name
	}

	if opts.Format == "" {
		opts.Format = "bullets"
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input Markdown file")
	}

	opts.Input = fs.Arg(0)

	return opts, nil
}
