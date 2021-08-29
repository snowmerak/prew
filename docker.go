package main

import "strings"

func convertPipPakcageToList(p []PipPackage) []PythonPackage {
	var packages []PythonPackage = nil
	queue := make([]PipPackage, len(p))
	copy(queue, p)
	cache := map[string]bool{}
	for len(queue) > 0 {
		pk := queue[0]
		queue = queue[1:]
		packages = append(packages, PythonPackage{
			Name:    pk.PackageName,
			Version: pk.InstalledVersion,
		})
		queue = append(queue, pk.Dependencies...)
	}
	for i := 0; i < len(packages)/2; i++ {
		packages[i], packages[len(packages)-1-i] = packages[len(packages)-1-i], packages[i]
	}
	for i := 0; i < len(packages); i++ {
		if cache[packages[i].Name] {
			packages = append(packages[:i], packages[i+1:]...)
			i--
			continue
		}
		cache[packages[i].Name] = true
	}
	return packages
}

func convertSpecToDockerfile(spec *Spec) string {
	sb := strings.Builder{}
	version := spec.Version
	if len(strings.Split(version, ".")) > 2 {
		sp := strings.Split(version, ".")
		version = sp[0] + "." + sp[1]
	}
	sb.WriteString("FROM python:" + version + "-buster\n\n")

	sb.WriteString("COPY src /src\n")

	ls := convertPipPakcageToList(spec.Dependencies)
	for _, v := range ls {
		sb.WriteString("RUN python3 -m pip install " + v.Name + "==" + v.Version + "\n")
	}

	sb.WriteString("WORKDIR /src")
	sb.WriteByte('\n')

	sb.WriteString("ENTRYPOINT python3 /src/" + spec.EntryFile + "\n")

	return sb.String()
}
