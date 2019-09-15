package main

import (
	"flag"
	"log"
	"os"

	"github.com/budougumi0617/dkl"
)

func main() {
	log.SetFlags(0)
	err := dkl.Run(os.Args, os.Stdin, os.Stdout, os.Stderr)
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		exitCode := 1
		if ecoder, ok := err.(interface{ ExitCode() int }); ok {
			exitCode = ecoder.ExitCode()
		}
		os.Exit(exitCode)
	}
}
