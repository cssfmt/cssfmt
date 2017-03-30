package main

import (
	"flag"
	"os"

	"github.com/gramework/gramework"
)

var (
	log      = gramework.Logger
	w        = flag.Bool("w", false, "rewrite files instead of printing fixed version")
	e        = flag.Bool("e", false, "show all errors, not the first 10 lines")
	d        = flag.Bool("d", false, "show diff instead of rewriting files")
	v        = flag.Bool("v", false, "validation mode")
	exitCode = 0
)

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		if err := processFile("<standard input>", os.Stdin, os.Stdout); err != nil {
			report("can't process standard input: %s", err)
		}
	}

	for i := 0; i < flag.NArg(); i++ {
		arg := flag.Arg(i)
		switch dir, err := os.Stat(arg); {
		case err != nil:
			report("can't process %q: %s", arg, err)
		case dir.IsDir():
			processDir(arg)
		default:
			if err := processFile(arg, nil, nil); err != nil {
				report("can't process %q: %s", arg, err)
			}
		}
	}

	os.Exit(exitCode)
}
