package main

import "strings"

func convertSpecToDockerfile(spec *Spec) string {
	sb := strings.Builder{}
	version := spec.Version
	if len(strings.Split(version, ".")) > 2 {
		sp := strings.Split(version, ".")
		version = sp[0] + "." + sp[1]
	}
	sb.WriteString("FROM python:" + version + "-buster\n\n")

	sb.WriteString("COPY src /src\n")

	for _, v := range spec.Dependencies {
		sb.WriteString("RUN python3 -m pip install " + v.Name + "==" + v.Version + "\n")
	}

	sb.WriteString("WORKDIR /src")
	sb.WriteByte('\n')

	sb.WriteString("ENTRYPOINT python3 /src/" + spec.EntryFile + "\n")

	return sb.String()
}
