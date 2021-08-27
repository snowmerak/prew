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

	remove := app.Command("remove", "Remove a package")

	run := app.Command("run", "Run python code")

	restore := app.Command("restore", "Restore a package")
	restorePath := restore.Arg("path", "Path to the package").Required().String()

	makeDockerfile := app.Command("make-dockerfile", "Create a Dockerfile")

	app.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case init.FullCommand():
		log.Println("Initializing project in current path")
		if err := initProjectToDirectory(*initPath); err != nil {
			log.Fatal(err)
		}
	case install.FullCommand():
		log.Println("Installing package", *installName)
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		spec, err := readSpecFromPath(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := appendDependencyToSpec(&spec, *installName); err != nil {
			log.Fatal(err)
		}
		if err := writeSpecToPath(path, spec); err != nil {
			log.Fatal(err)
		}
	case remove.FullCommand():
		log.Println("Removing package")
		path, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		spec, err := readSpecFromPath(path)
		if err != nil {
			log.Fatal(err)
		}
		if err := subductDependencyFromSpec(&spec); err != nil {
			log.Fatal(err)
		}
		if err := writeSpecToPath(path, spec); err != nil {
			log.Fatal(err)
		}
	case run.FullCommand():
		spec, err := readSpecFromPath(".")
		if err != nil {
			log.Fatal(err)
		}
		if err := runPythonCode(spec); err != nil {
			log.Fatal(err)
		}
	case restore.FullCommand():
		log.Println("Restore project")
		if err := restoreProject(*restorePath); err != nil {
			log.Fatal(err)
		}
	case makeDockerfile.FullCommand():
		log.Println("Create Dockerfile")
		spec, err := readSpecFromPath(".")
		if err != nil {
			log.Fatal(err)
		}
		data := convertSpecToDockerfile(&spec)
		if err := os.WriteFile("./dockerfile", []byte(data), 0644); err != nil {
			log.Fatal(err)
		}
	default:
		kingpin.Usage()
	}
}
