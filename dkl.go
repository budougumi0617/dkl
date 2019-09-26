package dkl

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os/exec"
)

const (
	version = "0.0.0"
)

// Run is entry point.
func Run(args []string, inStream io.Reader, outStream, errStream io.Writer) error {
	var v, d, k bool
	vdesc := "Print version information and quit."
	ddesc := "list docker containers and exec selected container."
	kdesc := "list pods and exec selected pod."
	cmdName := args[0]
	flags := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	flags.SetOutput(errStream)

	flags.BoolVar(&v, "version", false, vdesc)
	flags.BoolVar(&v, "v", false, vdesc)
	flags.BoolVar(&d, "docker", false, ddesc)
	flags.BoolVar(&d, "d", false, ddesc)
	flags.BoolVar(&k, "kubernetes", false, kdesc)
	flags.BoolVar(&k, "k", false, kdesc)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if v && d || v && k || d && k {
		m := "two or more options cannot be enabled."
		fmt.Fprintf(errStream, m)
		return fmt.Errorf(m)
	}

	nargs := flags.Args()
	if len(nargs) > 1 {
		m := "non-flag option must be one or less."
		fmt.Fprintf(errStream, m)
		return fmt.Errorf(m)
	}

	// バージョン情報の表示
	if v {
		fmt.Fprintf(errStream, "%s version %s\n", cmdName, version)
		return nil
	}

	cmd := []string{"/bin/bash"}
	var execCmd []string
	if len(nargs) > 0 {
		cmd[0] = nargs[0]
	}
	if k {
		pod, err := getPods()
		if err != nil {
			return err
		}
		execCmd = buildKubernetesCmd(pod, cmd)
	} else {
		container, err := show()
		if err != nil {
			return err
		}
		execCmd = buildDockerCmd(container.ID, cmd)
	}

	ecmd := exec.CommandContext(context.Background(), execCmd[0], execCmd[1:]...)
	ecmd.Stdin = inStream
	ecmd.Stdout = outStream
	ecmd.Stderr = errStream
	return ecmd.Run()
}
