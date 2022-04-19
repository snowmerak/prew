package main

import (
	"os"
	"os/exec"
)

// runExecCommand runs the given command and returns the output.
func runExecCommand(cmds []string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
