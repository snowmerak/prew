package main

import "runtime"

// python is the name of the python executable.
// when using windows, it is "python"
// when using another, it is "python3"
var python string

// pip3 is the name of the pip3 executable.
// when using windows, it is "Scripts\\pip3.exe"
// when using another, it is "bin/pip3"
var pip3 string

// vpython is the name of the virtualenv executable.
// when using windows, it is "Scripts\\python.exe"
// when using another, it is "bin/python3"
var vpython string

// virtualenv is the name of the virtualenv executable.
var virtualenv = "virtualenv==20.7.2"

func init() {
	if runtime.GOOS == "windows" {
		python = "python"
		pip3 = "Scripts\\pip3.exe"
		vpython = "Scripts\\python.exe"
	} else {
		python = "python3"
		pip3 = "bin/pip3"
		vpython = "bin/python3"
	}
}
