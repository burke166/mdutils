package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type target struct {
	goos   string
	goarch string
}

// Platforms built by --release. All dependencies are pure Go (no cgo), so
// these cross-compile cleanly from any host.
var releaseTargets = []target{
	{"windows", "amd64"},
	{"linux", "amd64"},
	{"linux", "arm64"}, // Raspberry Pi (64-bit OS) and AWS Graviton
	{"darwin", "amd64"},
	{"darwin", "arm64"},
}

func main() {
	release := flag.Bool("release", false, "build every cmd/* binary for all release targets into dist/<goos>_<goarch>/ instead of a single host build into bin/")
	flag.Parse()

	commands, err := commandDirs()
	if err != nil {
		fatal(err)
	}

	if *release {
		for _, t := range releaseTargets {
			outDir := filepath.Join("dist", t.goos+"_"+t.goarch)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				fatal(err)
			}

			for _, command := range commands {
				build(command, outDir, t.goos, t.goarch)
			}
		}
		return
	}

	if err := os.MkdirAll("bin", 0755); err != nil {
		fatal(err)
	}

	for _, command := range commands {
		build(command, "bin", "", "")
	}
}

func commandDirs() ([]string, error) {
	entries, err := filepath.Glob("cmd/*")
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		info, err := os.Stat(entry)
		if err != nil || !info.IsDir() {
			continue
		}

		if !hasMain(entry) {
			continue
		}

		dirs = append(dirs, entry)
	}

	return dirs, nil
}

// build compiles the command in cmdDir into outDir. An empty goos/goarch
// builds for the host platform using the ambient toolchain; otherwise it
// cross-compiles with CGO disabled.
func build(cmdDir, outDir, goos, goarch string) {
	name := filepath.Base(cmdDir)
	output := filepath.Join(outDir, executableName(name, goos))

	label := name
	if goos != "" {
		label = fmt.Sprintf("%s (%s/%s)", name, goos, goarch)
	}
	fmt.Printf("Building %s...", label)

	cmd := exec.Command(
		"go",
		"build",
		"-trimpath",
		"-ldflags=-s -w",
		"-o",
		output,
		"./"+cmdDir,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if goos != "" {
		cmd.Env = append(os.Environ(),
			"GOOS="+goos,
			"GOARCH="+goarch,
			"CGO_ENABLED=0",
		)
	}

	if err := cmd.Run(); err != nil {
		fatal(err)
	}

	fmt.Println("done!")
}

// executableName returns name with a .exe suffix on Windows targets. An
// empty targetOS means "the host platform".
func executableName(name, targetOS string) string {
	if targetOS == "" {
		targetOS = runtime.GOOS
	}

	if targetOS == "windows" {
		return name + ".exe"
	}

	return name
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func hasMain(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "main.go"))
	return err == nil
}
