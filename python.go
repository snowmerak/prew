package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

// runPythonCode runs the given python file of spec.yaml and returns the output.
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

// searchPackageVersion searches the given package version.
func searchPackageVersion(name string) error {
	version, err := selectPackageVersion(name)
	if err != nil {
		return err
	}
	reply := false
	if err := survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("Do you want to install %v version %v?", name, version)}, &reply, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if reply {
		spec, err := readSpecFromPath(".")
		if err != nil {
			return err
		}
		if err := appendDependencyToSpec(&spec, name, version); err != nil {
			return err
		}
		if err := writeSpecToPath(".", spec); err != nil {
			return err
		}
	}
	return nil
}

// getPythonVersion returns the python version.
func getPythonVersion() (string, error) {
	bs, err := exec.Command(python, "-V").Output()
	if err != nil {
		return "3.7.2", err
	}
	return strings.Split(strings.Trim(string(bs), "\r\n"), " ")[1], nil
}

// convertVersionToIntArray converts the given version to int array.
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

// compareVersion compares the given version.
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
