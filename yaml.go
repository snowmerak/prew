package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

// Spec is the specification of the project.
type Spec struct {
	Name         string       `yaml:"name"`
	Version      string       `yaml:"version"`
	EntryFile    string       `yaml:"entry_file"`
	Dependencies []PipPackage `yaml:"dependencies"`
}

// readSpecFromPath reads spec from path.
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

// writeSpecToPath writes spec to path.
func writeSpecToPath(path string, spec Spec) error {
	path = filepath.Join(path, "spec.yaml")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(spec)
}

// makeSpecFromSurvey makes spec from survey.
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
