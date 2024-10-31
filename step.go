package main

import "os/exec"

type step struct {
	name    string
	exe     string
	args    []string
	message string
	project string
}

func newStep(name string, exe string, args []string, message string, project string) step {
	return step{
		name,
		exe,
		args,
		message,
		project,
	}
}

func (x step) execute() (string, error) {
	cmd := exec.Command(x.exe, x.args...)
	cmd.Dir = x.project

	if err := cmd.Run(); err != nil {
		return "", &stepError{
			step:    x.name,
			message: "failed to execute",
			cause:   err,
		}
	}

	return x.message, nil
}
