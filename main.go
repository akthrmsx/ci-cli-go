package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type executor interface {
	execute() (string, error)
}

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

	pipeline := make([]executor, 4)
	pipeline[0] = newStep(
		"Go Build",
		"go",
		[]string{"build", ".", "errors"},
		"Go Build: SUCCESS",
		project,
	)
	pipeline[1] = newStep(
		"Go Test",
		"go",
		[]string{"test", "-v"},
		"Go Test: SUCCESS",
		project,
	)
	pipeline[2] = newExceptionStep(
		"Go Format",
		"gofmt",
		[]string{"-l", "."},
		"Go Format: SUCCESS",
		project,
	)
	pipeline[3] = newTimeoutStep(
		"Git Push",
		"git",
		[]string{"push", "origin", "main"},
		"Git Push: SUCCESS",
		project,
		10*time.Second,
	)

	signalChannel := make(chan os.Signal, 1)
	errChannel := make(chan error)
	doneChannel := make(chan struct{})

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for _, executor := range pipeline {
			message, err := executor.execute()

			if err != nil {
				errChannel <- err
				return
			}

			_, err = fmt.Fprintln(out, message)

			if err != nil {
				errChannel <- err
				return
			}
		}

		close(doneChannel)
	}()

	for {
		select {
		case received := <-signalChannel:
			signal.Stop(signalChannel)
			return fmt.Errorf("%s: Exiting: %w", received, ErrSignal)
		case err := <-errChannel:
			return err
		case <-doneChannel:
			return nil
		}
	}
}
