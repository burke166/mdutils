package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if err := os.MkdirAll("bin", 0755); err != nil {
		fatal(err)
	}

	commands, err := filepath.Glob("cmd/*")
	if err != nil {
		fatal(err)
	}

	for _, command := range commands {
		info, err := os.Stat(command)
		if err != nil || !info.IsDir() {
			continue
		}

		hasMain := hasMain(command)
		if err != nil {
			fatal(err)
		}
		if !hasMain {
			continue
		}

		name := filepath.Base(command)

		output := filepath.Join("bin", executableName(name))

		fmt.Printf("Building %s...", name)

		cmd := exec.Command(
			"go",
			"build",
			"-trimpath",
			"-ldflags=-s -w",
			"-o",
			output,
			"./"+command,
		)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fatal(err)
		}

		fmt.Println("done!")
	}
}

func executableName(name string) string {
	if os.PathSeparator == '\\' {
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
