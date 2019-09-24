package dkl

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getPods() (*v1.Pod, error) {
	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	pods, err := client.CoreV1().Pods("").List(meta_v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("pods.Item[0] = %+v\n", pods.Items[0])

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\u2388 {{ .Name | cyan }} ({{ .Namespace | red }} {{ .Status.Phase }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Namespace | red }} {{ .Status.Phase }})",
		Selected: "\u2388 executed {{ .Name | red }}",
		Details: `
--------- Pod ----------
{{ "Name:" | faint }}	{{ printf "%q" .Name }}
{{ "Namespacege:" | faint }}	{{ .Namespace }}
{{ "Status:" | faint }}	{{ .Status.Phase }}
{{ "Age:" | faint }}	{{ now | since  | print }}`,
		// {{ "Age:" | faint }}	{{ now | since .CreationTimestamp | print }}`,
		FuncMap: promptui.FuncMap,
	}
	templates.FuncMap["since"] = time.Since
	templates.FuncMap["now"] = time.Now

	searcher := func(input string, index int) bool {
		pod := pods.Items[index]
		ns := pod.Name + pod.Namespace
		lni := strings.ReplaceAll(strings.ToLower(ns), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(lni, input)
	}

	prompt := promptui.Select{
		Label:     "Kubernetes pod running...",
		Items:     pods.Items,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	return &pods.Items[i], nil
}

func newClient() (kubernetes.Interface, error) {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(kubeConfig)
}

func buildKubernetesCmd(pod *v1.Pod, cmd []string) []string {
	res := []string{"kubectl", "exec", "-it", pod.Name}
	res = append(res, cmd...)
	res = append(res, "-n", pod.Namespace)
	return res
}
