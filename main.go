package main

import (
	"flag"
	"fmt"
	"io"
	"os"
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
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}

	pipeline := make([]step, 1)
	pipeline[0] = newStep(
		"go build",
		"go",
		[]string{"build", ".", "errors"},
		"Go Build: SUCCESS",
		project,
	)

	for _, step := range pipeline {
		message, err := step.execute()

		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, message)

		if err != nil {
			return err
		}
	}

	return nil
}
