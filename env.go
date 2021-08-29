package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runExecCommand(cmds []string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installVirtualEnv() error {
	return runExecCommand([]string{python, "-m", "pip", "install", virtualenv})
}

func createVirtualEnv(path, version string) error {
	sp := strings.Split(version, ".")
	if len(sp) >= 3 {
		version = fmt.Sprintf("%s.%s", sp[0], sp[1])
	}
	return runExecCommand([]string{python, "-m", "virtualenv", path, "-p", fmt.Sprintf("python%s", version)})
}
