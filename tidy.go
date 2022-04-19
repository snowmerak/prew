package main

import (
	"log"

	"github.com/snowmerak/prew/pypi"
)

// tidyUpProject removes unused packages from the project.
func tidyUpProject(path string, yes bool) error {
	used, unused, expected, err := getUnusedPackages(path)
	if err != nil {
		return err
	}
	used["pipdeptree"] = true
	used["mypy"] = true
	dep := map[string]bool{}
	for k := range used {
		ls, err := getDependencies(path, k)
		if err != nil {
			return err
		}
		for _, l := range ls {
			dep[l] = true
		}
	}
	for k := range unused {
		if dep[k] {
			delete(unused, k)
		}
	}
	log.Println("Remove unused package")
	for k := range unused {
		if k == "pip" || k == "pipdeptree" || k == "mypy" {
			continue
		}
		log.Println("remove: ", k)
		if err := removePackage(k, yes); err != nil {
			return err
		}
	}
	log.Println("Install expected package")
	for k := range expected {
		_, err := pypi.GetPackageInfo(k, "")
		if err != nil {
			continue
		}
		log.Println("install: ", k)
		if err := installPackage(k, ""); err != nil {
			return err
		}
	}
	spec, err := readSpecFromPath(path)
	if err != nil {
		return err
	}
	packages, err := getDependencyTreeFrom(path)
	if err != nil {
		return err
	}
	spec.Dependencies = packages
	return writeSpecToPath(path, spec)
}
