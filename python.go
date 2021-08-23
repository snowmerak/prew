package main

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

func getPythonVersion() (string, error) {
	bs, err := exec.Command(python, "-V").Output()
	if err != nil {
		return "3.7.2", err
	}
	return strings.Split(strings.Trim(string(bs), "\n"), " ")[1], nil
}

func convertToPythonVersion(s string) [3]int {
	v := [3]int{0, 0, 0}
	sp := strings.Split(strings.TrimSpace(s), ".")
	var err error = nil
	for c, i := range sp {
		v[c], err = strconv.Atoi(i)
		if err != nil {
			return [3]int{0, 0, 0}
		}
	}
	return v
}

func checkValidVersion(v [3]int, s [3]int) bool {
	for c := range v {
		if v[c] < s[c] {
			return false
		}
	}
	return true
}

func installPackage(name string) error {
	_, err := exec.Command(python, "-m", "pip", "install", name).Output()
	if err != nil {
		return err
	}
	return nil
}

func installPackages(packages []string) error {
	for _, p := range packages {
		err := installPackage(p)
		if err != nil {
			return err
		}
	}
	return nil
}

type PythonPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func getInstalledPackageList() []PythonPackage {
	bs, err := exec.Command(python, "-m", "pip", "list", "--format=json", "--no-index").Output()
	if err != nil {
		return nil
	}
	var packages []PythonPackage
	err = json.Unmarshal(bs, &packages)
	if err != nil {
		return nil
	}
	return packages
}

func stringToPythonVersion(s string) [3]int {
	v := [3]int{0, 0, 0}
	sp := strings.Split(strings.TrimSpace(s), ".")
	var err error = nil
	for c, i := range sp {
		v[c], err = strconv.Atoi(i)
		if err != nil {
			return [3]int{0, 0, 0}
		}
	}
	return v
}

func removePackage(name string) error {
	_, err := exec.Command(python, "-m", "pip", "uninstall", name).Output()
	if err != nil {
		return err
	}
	return nil
}

func removePackages(packages []string) error {
	for _, p := range packages {
		err := removePackage(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func installPackageVersion(name string, version [3]int) error {
	err := installPackage(combinePackageNameVersion(name, version))
	if err != nil {
		return err
	}
	return nil
}

func combinePackageNameVersion(name string, version [3]int) string {
	return name + "==" + pythonVersionToString(version)
}

func pythonVersionToString(v [3]int) string {
	return strconv.Itoa(v[0]) + "." + strconv.Itoa(v[1]) + "." + strconv.Itoa(v[2])
}
