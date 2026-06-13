package mdstats

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/computercodeblue/mdutils/internal/fileutil"
)

type Options struct {
	Input       string
	Format      string
	NoRecursive bool
	Summary     bool
	PerFile     bool
	Exclude     []string
}

type excludeFlags []string

func (e *excludeFlags) String() string {
	return strings.Join(*e, ", ")
}

func (e *excludeFlags) Set(value string) error {
	*e = append(*e, value)
	return nil
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

	files, root, err := collectFiles(opts)
	if err != nil {
		return 3, err
	}

	files = fileutil.FilterExcluded(files, root, opts.Exclude)
	if len(files) == 0 {
		return 0, nil
	}

	var fileStats []FileStats
	readErrors := 0

	for _, file := range files {
		stats, err := AnalyzePath(file)
		if err != nil {
			fmt.Fprintf(stderr, "error: %s: %v\n", file, err)
			readErrors++
			continue
		}
		fileStats = append(fileStats, stats)
	}

	if len(fileStats) == 0 && readErrors > 0 {
		return 2, nil
	}

	result := AnalysisResult{
		Files:   fileStats,
		Summary: Summarize(fileStats),
	}

	output, err := renderOutput(result, opts)
	if err != nil {
		return 3, err
	}

	if output != "" {
		if _, err := fmt.Fprint(stdout, output); err != nil {
			return 3, err
		}
	}

	if readErrors > 0 {
		return 2, nil
	}

	return 0, nil
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

func renderOutput(result AnalysisResult, opts Options) (string, error) {
	switch opts.Format {
	case "json":
		return RenderJSON(result)
	case "csv":
		return RenderCSV(result.Files)
	default:
		return RenderText(result, opts.PerFile, opts.Summary), nil
	}
}

func parseOptions(args []string, stderr io.Writer) (Options, error) {
	var opts Options
	var excludes excludeFlags

	fs := flag.NewFlagSet("mdstats", flag.ContinueOnError)
	fs.SetOutput(stderr)

	jsonOut := fs.Bool("json", false, "output JSON")
	csvOut := fs.Bool("csv", false, "output CSV")
	fs.BoolVar(&opts.NoRecursive, "no-recursive", false, "do not analyze subdirectories")
	fs.BoolVar(&opts.Summary, "summary", false, "print a combined summary for all files")
	fs.BoolVar(&opts.PerFile, "per-file", true, "print statistics for each file")
	fs.Var(&excludes, "exclude", "exclude files matching a glob pattern")

	fs.Usage = func() {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "usage: %s [--json|--csv] [--summary] [--per-file] [--no-recursive] [--exclude glob] file-or-directory\n", exe)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return opts, err
	}

	formats := map[string]bool{
		"json": *jsonOut,
		"csv":  *csvOut,
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

	if fs.NArg() != 1 {
		fs.Usage()
		return opts, errors.New("missing input file or directory")
	}

	opts.Input = fs.Arg(0)
	opts.Exclude = excludes

	return opts, nil
}
