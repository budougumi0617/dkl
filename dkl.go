package dkl

import (
	"flag"
	"fmt"
	"io"
)

const (
	version = "0.0.0"
)

// Run is entry point.
func Run(args []string, inStream io.Reader, outStream, errStream io.Writer) error {
	var v bool
	cmdName := args[0]
	vdesc := "Print version information and quit."
	flags := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	flags.SetOutput(errStream)

	flags.BoolVar(&v, "version", false, vdesc)
	flags.BoolVar(&v, "v", false, vdesc)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// バージョン情報の表示
	if v {
		fmt.Fprintf(errStream, "%s version %s\n", cmdName, version)
		return nil
	}

	show()
	fmt.Fprint(outStream, "now implementing...\n")

	return nil
}
