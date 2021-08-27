package main

import "runtime"

var python string

var pip3 string
var vpython string
var activate string

func init() {
	if runtime.GOOS == "windows" {
		python = "python"
		pip3 = "Scripts\\pip3.exe"
		vpython = "Scripts\\python.exe"
		activate = "Scripts\\activate_this.py"
	} else {
		python = "python3"
		pip3 = "bin/pip3"
		vpython = "bin/python3"
		activate = "bin/activate_this.py"
	}
}
