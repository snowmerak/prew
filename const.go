package main

import "runtime"

var python string

var pip3 string
var vpython string

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
