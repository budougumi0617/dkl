package dkl

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

func show() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F40B {{ index .Names 0 | cyan }} ({{ .Image | red }}, {{ .ImageID }})",
		Inactive: "   {{ index .Names 0 | cyan }} ({{ .Image | red }}, {{ .ImageID  }})",
		Selected: "\U0001F336 {{ index .Names 0 | red }}",
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
		return
	}

	fmt.Printf("containers[%d] = %+v\n", i, containers[i])
}
