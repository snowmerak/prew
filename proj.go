package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func initProjectToDirectory(path string) error {
	spec := makeSpecFromSurvey()
	fmt.Println(spec)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if err := writeSpecToYaml(filepath.Join(path, "spec.yaml"), spec); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(path, "src"), 0755); err != nil {
		return err
	}
	app, err := os.Create(filepath.Join(path, "src", spec.EntryFile))
	if err != nil {
		return err
	}
	defer app.Close()
	if _, err := app.WriteString(`
	# Write your code here
	`); err != nil {
		return err
	}
	return nil
}
