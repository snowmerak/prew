package main

import (
	"os/exec"
	"strings"
)

func checkMypyInstalled(path string) bool {
	cmd := exec.Command(pip3, "list")
	data, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "mypy")
}

func installMypy(path string) error {
	if checkMypyInstalled(path) {
		return nil
	}
	_, err := exec.Command(pip3, "install", "mypy").Output()
	return err
}

func checkTypeFile(path, file string) ([]string, error) {
	cmd := exec.Command(vpython, "-m", "mypy", file)
	rs, _ := cmd.Output()
	return strings.Split(string(rs), "\n"), nil
}

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
