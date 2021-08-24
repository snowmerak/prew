package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

type Spec struct {
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version"`
	EntryFile    string       `yaml:"entry_file"`
	Dependencies []Dependency `yaml:"dependencies"`
}

type Dependency struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func readSpecFromCurrentPath() (Spec, error) {
	spec := Spec{}
	pw, err := os.Getwd()
	if err != nil {
		return spec, err
	}
	path := filepath.Join(pw, "spec.yaml")
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

func writeSpecToCurrentPath(spec Spec) error {
	pw, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(pw, "spec.yaml")
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
				Default: version,
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

func appendDependencyToSpec(spec *Spec, name, version string) error {
	if e := checkPackage(name, version); e != NotExist {
		if err := removePackage(name); err != nil {
			return err
		}
	}
	if err := installPackage(name, version); err != nil {
		return err
	}

	for i, d := range spec.Dependencies {
		if d.Name == name {
			switch compareVersion(d.Version, version) {
			case -1:
				fallthrough
			case 1:
				spec.Dependencies[i].Version = version
				return nil
			case 0:
				return nil
			}
		}
	}

	spec.Dependencies = append(spec.Dependencies, Dependency{Name: name, Version: version})
	return nil
}

func subductDependencyFromSpec(spec *Spec, name string) error {
	if e := checkPackage(name, ""); e != NotExist {
		if err := removePackage(name); err != nil {
			return err
		}
	}
	for i, d := range spec.Dependencies {
		if d.Name == name {
			spec.Dependencies = append(spec.Dependencies[:i], spec.Dependencies[i+1:]...)

			return nil
		}
	}
	return nil
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
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
