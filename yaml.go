package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"prew/pypi"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

type Spec struct {
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version"`
	EntryFile    string       `yaml:"entry_file"`
	Dependencies []PipPackage `yaml:"dependencies"`
}

func readSpecFromPath(path string) (Spec, error) {
	spec := Spec{}
	path = filepath.Join(path, "spec.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return spec, err
	}
	file, err := os.Open(path)
	if err != nil {
		return spec, err
	}
	defer file.Close()
	err = yaml.NewDecoder(file).Decode(&spec)
	return spec, err
}

func writeSpecToPath(path string, spec Spec) error {
	path = filepath.Join(path, "spec.yaml")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(spec)
}

func makeSpecFromSurvey() Spec {
	spec := Spec{}
	version, _ := getPythonVersion()
	ques := []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "project name: ",
			},
			Validate: survey.Required,
		},
		{
			Name: "Version",
			Prompt: &survey.Input{
				Message: "python version: ",
				Default: strings.TrimSpace(version),
			},
			Validate: survey.Required,
		},
		{
			Name: "EntryFile",
			Prompt: &survey.Input{
				Message: "entry file: ",
				Default: "app.py",
			},
			Validate: survey.Required,
		},
	}
	survey.Ask(ques, &spec)
	return spec
}

func selectPackageVersion(name string) (string, error) {
	version := ""
	versions := []string{}
	pack, err := pypi.GetPackageInfo(name, "")
	if err != nil {
		return version, err
	}
	for k := range pack.Releases {
		versions = append(versions, k)
	}
	prompt := &survey.Select{
		Message:  "select package version to install:",
		Options:  versions,
		PageSize: 13,
	}
	survey.AskOne(prompt, &version)
	return version, nil
}

func appendDependencyToSpec(spec *Spec, name string) error {
	version, err := selectPackageVersion(name)
	if err != nil {
		return err
	}
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if e := checkPackage(name, version); e != NotExist {
		if err := removePackage(name, true); err != nil {
			return err
		}
	}
	if err := installPackage(name, version); err != nil {
		return err
	}

	packages, err := getDependencyTreeFrom(path)
	if err != nil {
		return err
	}
	spec.Dependencies = packages
	return nil
}

func selectRemovePackages(spec *Spec) []string {
	packages := []string{}
	names := []string{}
	for _, d := range spec.Dependencies {
		names = append(names, fmt.Sprintf("%s == %s", d.PackageName, d.InstalledVersion))
	}
	prompt := &survey.MultiSelect{
		Message:  "select packages to remove:",
		Options:  names,
		PageSize: 13,
	}
	survey.AskOne(prompt, &packages)
	return packages
}

func subductDependencyFromSpec(spec *Spec, yes bool) error {
	selected := selectRemovePackages(spec)
	fmt.Println(selected)
	for _, v := range selected {
		v = strings.Split(v, " == ")[0]
		if e := checkPackage(v, ""); e != NotExist {
			fmt.Println(e)
			if err := removePackage(v, yes); err != nil {
				return err
			}
		}
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}
	packages, err := getDependencyTreeFrom(path)
	if err != nil {
		return err
	}
	spec.Dependencies = packages
	return nil
}

type PythonPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func getInstalledPackageList() []PythonPackage {
	bs, err := exec.Command(pip3, "list", "--format=json", "--no-index").Output()
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

const (
	ExistOlderVersion = 0
	ExistNewerVersion = 1
	ExistSameVersion  = 2
	NotExist          = 3
)

func checkPackage(name, version string) int {
	for _, p := range getInstalledPackageList() {
		fmt.Println(p.Name, name)
		if p.Name == name {
			if version == "" {
				return ExistSameVersion
			}
			switch compareVersion(p.Version, version) {
			case -1:
				return ExistNewerVersion
			case 0:
				return ExistSameVersion
			case 1:
				return ExistOlderVersion
			}
		}
	}
	return NotExist
}
