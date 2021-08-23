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

	app.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case init.FullCommand():
		log.Println("Initializing project in", initPath)
		if err := initProjectToDirectory(*initPath); err != nil {
			log.Fatal(err)
		}
	default:
	}

	kingpin.Usage()
}
