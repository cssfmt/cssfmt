package main

import (
	"os"
	"path/filepath"
)

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil && isCSSFile(f) {
		err = processFile(path, nil, nil)
	}
	if err != nil {
		report("error while processing %q: %s", path, err)
	}
	return nil
}

func processDir(p string) {
	filepath.Walk(p, visitFile)
}
