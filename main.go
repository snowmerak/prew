package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/snowmerak/prew/color"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("prew", "a python management tool")

	init := app.Command("init", "Initialize a new project")
	initPath := init.Arg("path", "Path to the project").Required().String()

	install := app.Command("install", "Install a package")
	installName := install.Arg("name", "Name of the package").Required().String()

	remove := app.Command("remove", "Remove a package")
	removeYes := remove.Flag("yes", "Remove without confirmation").Short('y').Bool()

	run := app.Command("run", "Run python code")

	restore := app.Command("restore", "Restore a package")
	restorePath := restore.Arg("path", "Path to the package").Required().String()

	makes := app.Command("make", "Make something")
	makesDockerfile := makes.Flag("dockerfile", "Make dockerfile").Short('d').Bool()

	tidy := app.Command("tidy", "Tidy up the project")
	tidyYes := tidy.Flag("yes", "Auto remove unused packages").Short('y').Bool()

	checkType := app.Command("check", "Check type of a python code")
	checkTypeAll := checkType.Flag("all", "Check all files").Short('a').Bool()
	checkTypeName := checkType.Arg("name", "Name of the file").String()

	app.Version("0.2.0")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case init.FullCommand():
		log.Println("Initializing project in current path")
		if err := initProjectToDirectory(*initPath); err != nil {
			log.Fatal(err)
		}
		log.Println(color.Green + "success" + color.Reset)
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
		log.Println(color.Green + "success" + color.Reset)
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
		if err := subductDependencyFromSpec(&spec, *removeYes); err != nil {
			log.Fatal(err)
		}
		if err := writeSpecToPath(path, spec); err != nil {
			log.Fatal(err)
		}
		log.Println(color.Green + "success" + color.Reset)
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
		log.Println(color.Green + "success" + color.Reset)
	case tidy.FullCommand():
		log.Println("Tidying project")
		if err := tidyUpProject(".", *tidyYes); err != nil {
			log.Fatal(err)
		}
		log.Println(color.Green + "success" + color.Reset)
	case makes.FullCommand():
		spec, err := readSpecFromPath(".")
		if err != nil {
			log.Fatal(err)
		}
		if *makesDockerfile {
			log.Println("Making Dockerfile")
			data := convertSpecToDockerfile(&spec)
			if err := os.WriteFile("./dockerfile", []byte(data), 0644); err != nil {
				log.Fatal(err)
			}
			log.Println(color.Green + "success" + color.Reset)
		}
	case checkType.FullCommand():
		// path, err := os.Getwd()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		path := "."
		path = filepath.Join(path, "src")
		if err := installMypy(path); err != nil {
			log.Fatal(err)
		}
		if *checkTypeAll {
			log.Println("Checking all files")
			em, err := checkTypeFiles(path)
			if err != nil {
				log.Fatal(err)
			}
			for _, el := range em {
				for _, e := range el {
					e = strings.ReplaceAll(e, "error", color.Red+"error"+color.Reset)
					e = strings.ReplaceAll(e, "Success", color.Green+"Success"+color.Reset)
					fmt.Println(e)
				}
			}
		}
		if *checkTypeName != "" {
			log.Println("Checking file", *checkTypeName)
			el, err := checkTypeFile(path, *checkTypeName)
			if err != nil {
				log.Fatal(err)
			}
			for _, e := range el {
				e = strings.ReplaceAll(e, "error", color.Red+"error"+color.Reset)
				e = strings.ReplaceAll(e, "Success", color.Green+"Success"+color.Reset)
				fmt.Println(e)
			}
		}
	default:
		kingpin.Usage()
	}
}
