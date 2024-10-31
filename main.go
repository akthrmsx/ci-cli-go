package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required")
	}

	args := []string{"build", ".", "errors"}
	cmd := exec.Command("go", args...)
	cmd.Dir = project

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'go build' failed: %w", err)
	}

	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	return err
}
