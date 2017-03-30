package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func processFile(path string, in io.Reader, out io.Writer) error {
	if in == nil {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	src, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	result, err := process(bytes.NewBuffer(src))
	if err != nil {
		return err
	}

	if *d {
		data, err := diff(src, result, path)
		if err != nil {
			return err
		}
		result = []byte(data)
	}
	if out == nil {
		out = os.Stdout
	}

	out.Write(result)

	return nil
}
