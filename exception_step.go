package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type exceptionStep struct {
	step
}

func newExceptionStep(name string, exe string, args []string, message string, project string) exceptionStep {
	return exceptionStep{
		step{
			name,
			exe,
			args,
			message,
			project,
		},
	}
}

func (x exceptionStep) execute() (string, error) {
	cmd := exec.Command(x.exe, x.args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = x.project

	if err := cmd.Run(); err != nil {
		return "", &stepError{
			step:    x.name,
			message: "failed to execute",
			cause:   err,
		}
	}

	if out.Len() > 0 {
		return "", &stepError{
			step:    x.name,
			message: fmt.Sprintf("invalid format: %s", out.String()),
			cause:   nil,
		}
	}

	return x.message, nil
}
