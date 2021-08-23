package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("prew", "a python management tool")

	init := app.Command("init", "Initialize a new project")
	initPath := init.Arg("path", "Path to the project").Required().String()

	install := app.Command("install", "Install a package")
	installName := install.Arg("name", "Name of the package").Required().String()
	installVersion := install.Arg("version", "Version of the package").String()

	remove := app.Command("remove", "Remove a package")
	removeName := remove.Arg("name", "Name of the package").Required().String()

	run := app.Command("run", "Run python code")

	app.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case init.FullCommand():
		log.Println("Initializing project in", *initPath)
		if err := initProjectToDirectory(*initPath); err != nil {
			log.Fatal(err)
		}
	case install.FullCommand():
		if *installVersion == "" {
			log.Println("Installing package", *installName)
		} else {
			log.Println("Installing package", *installName, "version", *installVersion)
		}
		spec, err := readSpecFromCurrentPath()
		if err != nil {
			log.Fatal(err)
		}
		if err := appendDependencyToSpec(&spec, *installName, *installVersion); err != nil {
			log.Fatal(err)
		}
		if err := writeSpecToCurrentPath(spec); err != nil {
			log.Fatal(err)
		}
	case remove.FullCommand():
		log.Println("Removing package", *removeName)
		spec, err := readSpecFromCurrentPath()
		if err != nil {
			log.Fatal(err)
		}
		if err := subductDependencyFromSpec(&spec, *removeName); err != nil {
			log.Fatal(err)
		}
		if err := writeSpecToCurrentPath(spec); err != nil {
			log.Fatal(err)
		}
	case run.FullCommand():
		spec, err := readSpecFromCurrentPath()
		if err != nil {
			log.Fatal(err)
		}
		if err := runPythonCode(spec); err != nil {
			log.Fatal(err)
		}
	default:
		kingpin.Usage()
	}
}
