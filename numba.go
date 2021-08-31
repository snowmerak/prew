package main

import "github.com/AlecAivazis/survey/v2"

func makeFunction() string {
	data := struct {
		Function string `survey:"function"`
	}{}
	ques := []*survey.Question{
		{
			Name: "function",
			Prompt: &survey.Input{
				Message: "Enter function name:",
			},
			Validate: survey.Required,
		},
	}
	survey.Ask(ques, &data)
	return ""
}
