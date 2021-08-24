package main

import "runtime"

var python string

var pip3 string
var vpython string

func init() {
	if runtime.GOOS == "windows" {
		python = "python"
		pip3 = "bin/pip3.exe"
		vpython = "bin/python.exe"
	} else {
		python = "python3"
		pip3 = "bin/pip3"
		vpython = "bin/python3"
	}
}
