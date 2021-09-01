package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/snowmerak/prew/pypi"
)

func getFileListRecursive(path string) ([]string, error) {
	q := []string{path}
	l := []string{}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		fs, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			if f.IsDir() {
				q = append(q, filepath.Join(p, f.Name()))
			} else {
				if strings.HasSuffix(f.Name(), ".py") {
					l = append(l, filepath.Join(p, f.Name()))
				}
			}
		}
	}
	return l, nil
}

func getImportedPackage(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := regexp.MustCompile(`\s+`)
	l := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "import") {
			s := sc.ReplaceAllString(line, " ")
			s = strings.ReplaceAll(s, ",", "")
			l = append(l, strings.Split(s, " ")[1:]...)
			continue
		}
		if strings.HasPrefix(line, "from") {
			s := sc.ReplaceAllString(strings.TrimSpace(line), " ")
			s = strings.SplitN(s, " ", 3)[1]
			l = append(l, s)
			continue
		}
	}
	return l, nil
}

func getUsedPackages(path string) ([]string, error) {
	ls, err := getFileListRecursive(path)
	if err != nil {
		return nil, err
	}
	c := map[string]bool{}
	l := []string{}
	for _, f := range ls {
		ps, err := getImportedPackage(f)
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			if c[p] {
				continue
			}
			c[p] = true
			l = append(l, p)
		}
	}
	return l, nil
}

func getTopLevelFromPackages(path string) (map[string]string, error) {
	spec, err := readSpecFromPath(path)
	if err != nil {
		return nil, err
	}
	version := spec.Version
	if len(strings.Split(version, ".")) > 2 {
		sp := strings.Split(version, ".")
		version = sp[0] + "." + sp[1]
	}
	if runtime.GOOS == "windows" {
		path = filepath.Join(path, "Lib", "site-packages")
	} else {
		path = filepath.Join(path, "lib", "python"+version, "site-packages")
	}
	fs, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	m := map[string]string{}
	for _, f := range fs {
		if f.IsDir() && strings.HasSuffix(f.Name(), "dist-info") {
			pname := strings.SplitN(f.Name(), "-", 2)[0]
			tf, err := os.Open(filepath.Join(path, f.Name(), "top_level.txt"))
			if err != nil {
				return nil, err
			}
			defer tf.Close()
			sc := bufio.NewScanner(tf)
			for sc.Scan() {
				line := sc.Text()
				line = strings.TrimSpace(line)
				m[line] = pname
			}
		}
	}
	return m, nil
}

func getUnusedPackages(path string) (map[string]bool, map[string]bool, map[string]bool, error) {
	used, err := getUsedPackages(filepath.Join(".", "src"))
	if err != nil {
		return nil, nil, nil, err
	}
	top, err := getTopLevelFromPackages(".")
	if err != nil {
		return nil, nil, nil, err
	}
	expected := map[string]bool{}
	usedList := map[string]bool{}
	for _, u := range used {
		if _, ok := top[u]; !ok {
			expected[top[u]] = true
		} else {
			usedList[top[u]] = true
		}
		delete(top, u)
	}
	c := map[string]bool{}
	unusedList := map[string]bool{}
	for _, v := range top {
		if c[v] {
			continue
		}
		c[v] = true
		unusedList[v] = true
	}
	return usedList, unusedList, expected, nil
}

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
