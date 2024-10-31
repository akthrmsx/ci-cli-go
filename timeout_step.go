package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

func newTimeoutStep(
	name string,
	exe string,
	args []string,
	message string,
	project string,
	timeout time.Duration,
) timeoutStep {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return timeoutStep{
		step: step{
			name,
			exe,
			args,
			message,
			project,
		},
		timeout: timeout,
	}
}

func (x timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), x.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, x.exe, x.args...)
	cmd.Dir = x.project

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &stepError{
				step:    x.name,
				message: "failed time out",
				cause:   context.DeadlineExceeded,
			}
		}

		return "", &stepError{
			step:    x.name,
			message: "failed to execute",
			cause:   err,
		}
	}

	return x.message, nil
}
