package main

import (
	"fmt"
	"strings"
)

// installVirtualEnv installs virtualenv.
func installVirtualEnv() error {
	return runExecCommand([]string{python, "-m", "pip", "install", virtualenv})
}

// createVirtualEnv creates virtualenv.
func createVirtualEnv(path, version string) error {
	sp := strings.Split(version, ".")
	if len(sp) >= 3 {
		version = fmt.Sprintf("%s.%s", sp[0], sp[1])
	}
	return runExecCommand([]string{python, "-m", "virtualenv", path, "-p", fmt.Sprintf("python%s", version)})
}
