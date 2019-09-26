package dkl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/manifoldco/promptui"
)

func show() (*types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("running container is not found.")
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F40B {{ index .Names 0 | cyan }} ({{ .Image | red }}, {{ .ImageID }})",
		Inactive: "   {{ index .Names 0 | cyan }} ({{ .Image | red }}, {{ .ImageID  }})",
		Selected: "\U0001F40B {{ index .Names 0 | red }}",
		Details: `
--------- Image ----------
{{ "Name:" | faint }}	{{ printf "%q" .Names }}
{{ "Image:" | faint }}	{{ .Image }}
{{ "ImageID:" | faint }}	{{ .ImageID }}`,
	}

	searcher := func(input string, index int) bool {
		container := containers[index]
		n := strings.Join(container.Names, "")
		ni := n + container.Image
		lni := strings.ReplaceAll(strings.ToLower(ni), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(lni, input)
	}

	prompt := promptui.Select{
		Label:     "Docker running...",
		Items:     containers,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	return &containers[i], nil
}

// Port from moby/moby/integration/internal/container/exec.go

// ExecResult represents a result returned from Exec()
type ExecResult struct {
	ExitCode  int
	outBuffer *bytes.Buffer
	errBuffer *bytes.Buffer
}

// Stdout returns stdout output of a command run by Exec()
func (res *ExecResult) Stdout() string {
	return res.outBuffer.String()
}

// Stderr returns stderr output of a command run by Exec()
func (res *ExecResult) Stderr() string {
	return res.errBuffer.String()
}

// Combined returns combined stdout and stderr output of a command run by Exec()
func (res *ExecResult) Combined() string {
	return res.outBuffer.String() + res.errBuffer.String()
}

func execDocker(cid string, cmd []string, ins io.Reader) (ExecResult, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	econf := types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}
	cresp, err := cli.ContainerExecCreate(ctx, cid, econf)
	if err != nil {
		return ExecResult{}, err
	}
	execID := cresp.ID

	// run it, with stdout/stderr attached
	aresp, err := cli.ContainerExecAttach(ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}
	defer aresp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, aresp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return ExecResult{}, err
		}
		break

	case <-ctx.Done():
		return ExecResult{}, ctx.Err()
	}

	// get the exit code
	iresp, err := cli.ContainerExecInspect(ctx, execID)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{ExitCode: iresp.ExitCode, outBuffer: &outBuf, errBuffer: &errBuf}, nil
}

func buildDockerCmd(cid string, cmd []string) []string {
	res := []string{"docker", "exec", "-i", "-t", cid}
	res = append(res, cmd...)
	return res

}
