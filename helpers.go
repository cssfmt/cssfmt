package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	delta "github.com/octavore/delta/lib"
	deltafmt "github.com/octavore/delta/lib/formatter"
)

var reportedErrors = 0

func report(format string, v ...interface{}) {
	exitCode = 2
	if !*e && reportedErrors < 10 {
		log.Errorf(format, v...)
	}
	reportedErrors++
}

func diff(b1, b2 []byte, filename string) (data string, err error) {
	if runtime.GOOS != "windows" { // we have native `diff` in all operation systems.
		return nativediff(b1, b2, filename)
	}
	d := delta.HistogramDiff(string(b1), string(b2))
	data = deltafmt.Text(d)
	return
}

func isCSSFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() &&
		!strings.HasPrefix(name, ".") &&
		(strings.HasSuffix(name, ".css") || strings.HasSuffix(name, ".scss") ||
			strings.HasSuffix(name, ".sass") || strings.HasSuffix(name, ".less"))
}

func usage() {
	fmt.Printf(`Usage of %s: [flags] [path ...]
	-w write to source files, not to stdout
	-d display diffs instead of rewriting files
	-e show all errors, not the first 10 lines
	-v validate file silently and exit with code 1 if modification occurred
`, os.Args[0])
}

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the goimports file in
// 3rd-party licenses directory.
func nativediff(b1, b2 []byte, filename string) (data string, err error) {
	f1, err := ioutil.TempFile("", "cssfmt")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "cssfmt")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	cmd := "diff"
	if runtime.GOOS == "plan9" {
		cmd = "/bin/ape/diff"
	}

	out, err := exec.Command(cmd, "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(out) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil

		out, err = replaceTempFilename(out, filename)
	}
	return string(out), err
}

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the goimports file in
// 3rd-party licenses directory.
func replaceTempFilename(diff []byte, filename string) ([]byte, error) {
	bs := bytes.SplitN(diff, []byte{'\n'}, 3)
	if len(bs) < 3 {
		return nil, fmt.Errorf("got unexpected diff for %s", filename)
	}
	// Preserve timestamps.
	var t0, t1 []byte
	if i := bytes.LastIndexByte(bs[0], '\t'); i != -1 {
		t0 = bs[0][i:]
	}
	if i := bytes.LastIndexByte(bs[1], '\t'); i != -1 {
		t1 = bs[1][i:]
	}
	// Always print filepath with slash separator.
	f := filepath.ToSlash(filename)
	bs[0] = []byte(fmt.Sprintf("--- %s%s", f+".orig", t0))
	bs[1] = []byte(fmt.Sprintf("+++ %s%s", f, t1))
	return bytes.Join(bs, []byte{'\n'}), nil
}
