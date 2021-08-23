package main

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v2"
)

type Spec struct {
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	EntryFile    string `yaml:"entry_file"`
	Dependencies []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"dependencies"`
}

func readSpecFromYaml(path string) (Spec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, err
	}
	var spec Spec
	yaml.Unmarshal(data, &spec)
	return spec, err
}

func writeSpecToYaml(path string, spec Spec) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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
