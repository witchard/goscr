package main

import (
	"errors"
	"fmt"
	"os/exec"
)

func Run(dir string, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	command.Dir = dir
	out, err := command.Output()
	if err != nil {
		var execErr *exec.ExitError
		if errors.As(err, &execErr) {
			fmt.Println("Command", cmd, args, "exited with status", execErr.ProcessState.ExitCode())
			fmt.Println("---------- stdout ----------")
			fmt.Println(string(out))
			fmt.Println("---------- stdout ----------")
			fmt.Println(string(execErr.Stderr))
		}
	}
	return err
}

func Compile(workdir string) error {
	if err := Run(workdir, "go", "mod", "init", "goscr"); err != nil {
		return err
	}
	if err := Run(workdir, "goimports", "-w", "main.go"); err != nil {
		return err
	}
	if err := Run(workdir, "go", "mod", "tidy"); err != nil {
		return err
	}
	if err := Run(workdir, "go", "build"); err != nil {
		return err
	}
	return nil
}
