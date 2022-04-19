package main

import (
	"os"
	"path/filepath"
	"strings"
)

// getFileListRecursive returns the list of files in the given directory.
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
