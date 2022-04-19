package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
)

// getDependencyTreeFrom returns the dependency tree of the current directory.
func getDependencyTreeFrom(_ string) ([]PipPackage, error) {
	if checkPackage("pipdeptree", "") == NotExist {
		if err := installPackage("pipdeptree", ""); err != nil {
			return nil, err
		}
	}
	data, err := exec.Command(vpython, "-m", "pipdeptree", "--json-tree").Output()
	if err != nil {
		return nil, err
	}
	p := new([]PipPackage)
	*p = []PipPackage{}
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}
	return *p, nil
}

// getDependencies returns the dependencies of the current directory.
func getDependencies(path string, name string) ([]string, error) {
	spec, err := readSpecFromPath(path)
	if err != nil {
		return nil, err
	}
	c := map[string]bool{
		name: true,
	}
	r := map[string]bool{}
	al := make([]PipPackage, len(spec.Dependencies))
	copy(al, spec.Dependencies)

	for len(al) > 0 {
		p := al[0]
		al = al[1:]
		if c[p.PackageName] {
			for _, d := range p.Dependencies {
				if !c[d.PackageName] {
					c[d.PackageName] = true
					if strings.Contains(d.PackageName, "_") {
						c[strings.ReplaceAll(d.PackageName, "_", "-")] = true
					}
					if strings.Contains(d.PackageName, "-") {
						c[strings.ReplaceAll(d.PackageName, "-", "_")] = true
					}
				}
			}
		} else {
			for _, d := range p.Dependencies {
				r[d.PackageName] = true
				if strings.Contains(d.PackageName, "_") {
					r[strings.ReplaceAll(d.PackageName, "_", "-")] = true
				}
				if strings.Contains(d.PackageName, "-") {
					r[strings.ReplaceAll(d.PackageName, "-", "_")] = true
				}
			}
		}
		al = append(al, p.Dependencies...)
	}

	for k := range r {
		delete(c, k)
	}

	deps := make([]string, 0, len(c))
	for k := range c {
		deps = append(deps, k)
	}

	return deps, nil
}

// addDependencyToSpec adds dependencies to the current directory's spec.yaml.
func appendDependencyToSpec(spec *Spec, name, version string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	if e := checkPackage(name, version); e != NotExist {
		if err := removePackage(name, true); err != nil {
			return err
		}
	}
	if err := installPackage(name, version); err != nil {
		return err
	}

	packages, err := getDependencyTreeFrom(path)
	if err != nil {
		return err
	}
	spec.Dependencies = packages
	return nil
}

// subductDependencyFromSpec removes dependencies from the current directory's spec.yaml.
func subductDependencyFromSpec(spec *Spec, name string, yes bool, dep bool) error {
	target := []string{name}

	if dep {
		deps, err := getDependencies(".", name)
		if err != nil {
			return err
		}
		target = append(target, deps...)
	}

	for _, t := range target {
		if e := checkPackage(t, ""); e != NotExist {
			if err := removePackage(t, yes); err != nil {
				return err
			}
		}
	}

	packages, err := getDependencyTreeFrom(".")
	if err != nil {
		return err
	}
	spec.Dependencies = packages
	return nil
}
