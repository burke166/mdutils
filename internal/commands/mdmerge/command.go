package mdmerge

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Options struct {
	Inputs []string
	Out    string
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

	files, err := CollectMarkdownFiles(opts.Inputs)
	if err != nil {
		if IsNotExist(err) {
			return 2, err
		}
		fmt.Fprintln(stderr, "error:", err)
		return 1, nil
	}

	contents, err := ReadMarkdownFiles(files)
	if err != nil {
		return 2, err
	}

	merged := MergeContents(contents)

	if opts.Out == "" {
		if _, err := fmt.Fprint(stdout, merged); err != nil {
			return 2, err
		}
		return 0, nil
	}

	if err := os.WriteFile(opts.Out, []byte(merged), 0644); err != nil {
		return 2, err
	}

	return 0, nil
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	fs := flag.NewFlagSet("mdmerge", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&opts.Out, "out", "", "output Markdown file")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--out file] file-or-directory...\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if fs.NArg() == 0 {
		fs.Usage()
		return opts, errors.New("missing input files or directory")
	}

	opts.Inputs = fs.Args()

	return opts, nil
}
