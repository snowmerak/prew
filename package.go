package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/snowmerak/prew/pypi"
)

// PythonPackage is a struct of python package
type PythonPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// PipPackage is a struct of pip package
type PipPackage struct {
	PackageName      string       `json:"package_name" yaml:"name"`
	InstalledVersion string       `json:"installed_version" yaml:"version"`
	Dependencies     []PipPackage `json:"dependencies" yaml:"dependencies"`
}

// VersionCompared result of comparation of version.
const (
	ExistOlderVersion = 0
	ExistNewerVersion = 1
	ExistSameVersion  = 2
	NotExist          = 3
)

// printPackages prints packages
func printPackages() error {
	spec, err := readSpecFromPath(".")
	if err != nil {
		return err
	}
	packages := []string{}
	{
		pippackages := convertPipPakcageToList(spec.Dependencies)
		for _, p := range pippackages {
			packages = append(packages, p.Name)
		}
	}
	selected := ""
	survey.AskOne(&survey.Select{Message: "Select packages", Options: packages}, &selected, survey.WithPageSize(5))
	reply := false
	if err := survey.AskOne(&survey.Confirm{Message: fmt.Sprintf("Do you want to remove %v?", selected)}, &reply, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	if reply {
		if err := subductDependencyFromSpec(&spec, selected, false, false); err != nil {
			return err
		}
		if err := writeSpecToPath(".", spec); err != nil {
			return err
		}
	}
	return nil
}

// installPackage installs package
func installPackage(name, version string) error {
	cmd := exec.Command(pip3)
	if version != "" {
		cmd.Args = append(cmd.Args, "install", fmt.Sprintf("%v==%v", name, version))
	} else {
		cmd.Args = append(cmd.Args, "install", name)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Println(name, version, "installed")
	return nil
}

// removePackage removes package
func removePackage(name string, yes bool) error {
	cmd := exec.Command(pip3, "uninstall", name)
	if yes {
		cmd.Args = append(cmd.Args, "-y")
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	log.Println(name, "uninstalled")
	return nil
}

// getImportedPackage gets imported package of given file.
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

// getUsedPackages gets used packages of given path.
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

// getTopLevelFromPackage gets top level package name of given path's spec.yaml's packages.
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

// getUnusedPackages gets unused packages of root path.
func getUnusedPackages(_ string) (map[string]bool, map[string]bool, map[string]bool, error) {
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

// checkPackage check and compare package version.
func checkPackage(name, version string) int {
	for _, p := range getInstalledPackageList() {
		if p.Name == name {
			if version == "" {
				return ExistSameVersion
			}
			switch compareVersion(p.Version, version) {
			case -1:
				return ExistNewerVersion
			case 0:
				return ExistSameVersion
			case 1:
				return ExistOlderVersion
			}
		}
	}
	return NotExist
}

// getInstalledPackageList gets installed package list.
func getInstalledPackageList() []PythonPackage {
	bs, err := exec.Command(pip3, "list", "--format=json", "--no-index").Output()
	if err != nil {
		return nil
	}
	var packages []PythonPackage
	err = json.Unmarshal(bs, &packages)
	if err != nil {
		return nil
	}
	return packages
}

// selectPackagesVersion selects packages version.
func selectPackageVersion(name string) (string, error) {
	version := ""
	versions := []string{}
	pack, err := pypi.GetPackageInfo(name, "")
	if err != nil {
		return version, err
	}
	for k := range pack.Releases {
		versions = append(versions, k)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})
	prompt := &survey.Select{
		Message:  "package versions:",
		Options:  versions,
		PageSize: 5,
	}
	survey.AskOne(prompt, &version)
	return version, nil
}
