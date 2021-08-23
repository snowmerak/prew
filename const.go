package main

import "runtime"

var python string

func init() {
	if runtime.GOOS == "windows" {
		python = "python"
	} else {
		python = "python3"
	}
}
