package main

import (
	"os/exec"
	"strings"
)

// checkMypyInstalled returns true if mypy is installed.
func checkMypyInstalled(_ string) bool {
	cmd := exec.Command(pip3, "list")
	data, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "mypy")
}

// installMypy installs mypy.
func installMypy(path string) error {
	if checkMypyInstalled(path) {
		return nil
	}
	_, err := exec.Command(pip3, "install", "mypy==0.910").Output()
	return err
}

// checkTypeFile checks type of given file.
func checkTypeFile(_, file string) ([]string, error) {
	cmd := exec.Command(vpython, "-m", "mypy", file)
	rs, _ := cmd.Output()
	return strings.Split(string(rs), "\n"), nil
}

// checkTypeFiles checks type of given files.
func checkTypeFiles(path string) (map[string][]string, error) {
	fs, err := getFileListRecursive(path)
	if err != nil {
		return nil, err
	}
	rs := make(map[string][]string)
	for _, f := range fs {
		el, err := checkTypeFile(path, f)
		if err != nil {
			return nil, err
		}
		rs[f] = el
	}
	return rs, nil
}
