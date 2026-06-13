package mdlint

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/computercodeblue/mdutils/internal/fileutil"
)

type Options struct {
	Input       string
	ConfigPath  string
	JSON        bool
	Quiet       bool
	NoRecursive bool
}

func Execute() {
	code, err := Run(os.Args[1:], os.Stdout, os.Stderr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(3)
	}
	os.Exit(code)
}

func Run(args []string, stdout io.Writer, stderr io.Writer) (int, error) {
	opts, err := parseOptions(args, stderr)
	if err != nil {
		return 3, err
	}

	cfg, err := LoadConfig(opts.ConfigPath)
	if err != nil {
		return 3, err
	}

	files, root, err := collectFiles(opts)
	if err != nil {
		return 3, err
	}

	files = fileutil.FilterExcluded(files, root, cfg.Exclude)
	if len(files) == 0 {
		return 0, nil
	}

	issues, err := LintPaths(files, cfg)
	if err != nil {
		return 3, err
	}

	if len(issues) == 0 {
		return 0, nil
	}

	output, err := renderOutput(issues, opts.JSON, opts.Quiet)
	if err != nil {
		return 3, err
	}

	if output != "" {
		if _, err := fmt.Fprint(stdout, output); err != nil {
			return 3, err
		}
	}

	return exitCode(issues), nil
}

func collectFiles(opts Options) ([]string, string, error) {
	recursive := !opts.NoRecursive
	files, err := fileutil.CollectMarkdownFiles(opts.Input, recursive)
	if err != nil {
		return nil, "", err
	}

	root := ""
	if info, statErr := os.Stat(opts.Input); statErr == nil && info.IsDir() {
		root = opts.Input
	}

	return files, root, nil
}

func renderOutput(issues []Issue, asJSON bool, quiet bool) (string, error) {
	if asJSON {
		return RenderJSON(issues, quiet)
	}
	return RenderText(issues, quiet), nil
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options

	fs := flag.NewFlagSet("mdlint", flag.ContinueOnError)
	fs.SetOutput(stderr)

	fs.StringVar(&opts.ConfigPath, "config", "", "path to a YAML config file")
	fs.BoolVar(&opts.JSON, "json", false, "output issues as JSON")
	fs.BoolVar(&opts.Quiet, "quiet", false, "only print errors")
	fs.BoolVar(&opts.NoRecursive, "no-recursive", false, "do not lint subdirectories")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--config path] [--json] [--quiet] [--no-recursive] file-or-directory\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input file or directory")
	}

	opts.Input = fs.Arg(0)

	return opts, nil
}
