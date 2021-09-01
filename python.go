package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runPythonCode(spec Spec) error {
	for _, d := range spec.Dependencies {
		if e := checkPackage(d.PackageName, d.InstalledVersion); e == ExistSameVersion {
			continue
		}
		if err := installPackage(d.PackageName, d.InstalledVersion); err != nil {
			return err
		}
	}
	return runExecCommand([]string{vpython, filepath.Join("src", spec.EntryFile)})
}

func getPythonVersion() (string, error) {
	bs, err := exec.Command(python, "-V").Output()
	if err != nil {
		return "3.7.2", err
	}
	return strings.Split(strings.Trim(string(bs), "\r\n"), " ")[1], nil
}

func installPackage(name, version string) error {
	cmd := exec.Command(pip3)
	if version != "" {
		cmd.Args = append(cmd.Args, "install", fmt.Sprintf("%v==%v", name, version))
	} else {
		cmd.Args = append(cmd.Args, "install", name)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println(name, version, "installed")
	return nil
}

func removePackage(name string, yes bool) error {
	cmd := exec.Command(pip3, "uninstall", name)
	if yes {
		cmd.Args = append(cmd.Args, "-y")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	log.Println(name, "uninstalled")
	return nil
}

func convertVersionToIntArray(version string) [3]int {
	v := [3]int{0, 0, 0}
	s := strings.Split(strings.TrimSpace(version), ".")
	var err error = nil
	for i := 0; i < 3; i++ {
		if len(s) <= i {
			v[i] = 0
			continue
		}
		v[i], err = strconv.Atoi(s[i])
		if err != nil {
			return v
		}
	}
	return v
}

func compareVersion(a, b string) int {
	av := convertVersionToIntArray(a)
	bv := convertVersionToIntArray(b)
	for i := 0; i < 3; i++ {
		if av[i] > bv[i] {
			return -1
		} else if av[i] < bv[i] {
			return 1
		}
	}
	return 0
}

type PipPackage struct {
	PackageName      string       `json:"package_name" yaml:"name"`
	InstalledVersion string       `json:"installed_version" yaml:"version"`
	Dependencies     []PipPackage `json:"dependencies" yaml:"dependencies"`
}

func getDependencyTreeFrom(path string) ([]PipPackage, error) {
	if checkPackage("pipdeptree", "") == NotExist {
		if err := installPackage("pipdeptree", ""); err != nil {
			return nil, err
		}
	}
	data, err := exec.Command(vpython, "-m", "pipdeptree", "--json-tree").Output()
	if err != nil {
		return nil, err
	}
	p := new([]PipPackage)
	*p = []PipPackage{}
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return *p, nil
}
