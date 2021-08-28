package main

import (
	"os"
	"path/filepath"
	"strings"
)

func initProjectToDirectory(path string) error {
	spec := makeSpecFromSurvey()
	if err := installVirtualEnv(); err != nil {
		return err
	}
	if err := createVirtualEnv(path, spec.Version); err != nil {
		return err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	if err := writeSpecToPath(path, spec); err != nil {
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
	ignore, err := os.Create(filepath.Join(path, ".gitignore"))
	if err != nil {
		return err
	}
	defer ignore.Close()
	if _, err := ignore.WriteString(strings.TrimSpace("bin\nlib\npyvenv.cfg\ndockerfile")); err != nil {
		return err
	}
	return nil
}

func restoreProject(path string) error {
	spec, err := readSpecFromPath(path)
	if err != nil {
		return err
	}
	if err := createVirtualEnv(path, spec.Version); err != nil {
		return err
	}
	ls := convertPipPakcageToList(spec.Dependencies)
	for _, d := range ls {
		if e := checkPackage(d.Name, d.Version); e == ExistSameVersion {
			continue
		}
		if err := installPackage(d.Name, d.Version); err != nil {
			return err
		}
	}
	return nil
}
