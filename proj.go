package main

import (
	"os"
	"path/filepath"
)

func initProjectToDirectory() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	spec := makeSpecFromSurvey()
	if err := installVirtualEnv(); err != nil {
		return err
	}
	if err := createVirtualEnv(spec.Version); err != nil {
		return err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if err := writeSpecToCurrentPath(spec); err != nil {
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

func restoreProject() error {
	spec, err := readSpecFromCurrentPath()
	if err != nil {
		return err
	}
	if err := createVirtualEnv(spec.Version); err != nil {
		return err
	}
	for _, d := range spec.Dependencies {
		if e := checkPackage(d.Name, d.Version); e == ExistSameVersion {
			continue
		}
		if err := installPackage(d.Name, d.Version); err != nil {
			return err
		}
	}
	return nil
}
