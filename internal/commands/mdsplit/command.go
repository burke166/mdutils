package mdsplit

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Options struct {
	Input string
	Level int
	Out   string
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

	sections, err := SplitMarkdown(string(source), opts.Level)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1, nil
	}

	outDir := opts.Out
	if outDir == "" {
		outDir = "."
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		return 2, err
	}

	used := make(map[string]int)
	for _, section := range sections {
		filename := sectionFilename(section, used)
		path := filepath.Join(outDir, filename)
		if err := os.WriteFile(path, []byte(section.Content), 0644); err != nil {
			return 2, err
		}
	}

	return 0, nil
}

func sectionFilename(section Section, used map[string]int) string {
	if section.Heading == "" {
		return "00-preamble.md"
	}

	base := EnsureUniqueFilename(section.Slug, used)
	return base + ".md"
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options
	opts.Level = 1

	fs := flag.NewFlagSet("mdsplit", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.IntVar(&opts.Level, "level", 1, "heading level to split on")
	fs.StringVar(&opts.Out, "out", "", "destination folder")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--level N] [--out dir] file.md\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if err := validateLevel(opts.Level); err != nil {
		return opts, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input Markdown file")
	}

	opts.Input = fs.Arg(0)

	return opts, nil
}

func validateLevel(level int) error {
	if level < 1 {
		return errors.New("level must be at least 1")
	}
	if level > 6 {
		return errors.New("level must be at most 6")
	}
	return nil
}
