package main

import (
	"strings"

	"golang.org/x/net/html"
)

type Step struct {
	Element, Id, Class string
}

func step(n *html.Node, i int, steps []*Step, endFn func(*html.Node)) {
	if i == len(steps) {
		endFn(n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if IsValid(c, steps[i]) {
			step(c, i + 1, steps, endFn)
		}
	}
}

func IsValid(n *html.Node, step *Step) bool {
	if n.Type == html.ElementNode && n.Data == step.Element {
		id := GetAttr(n, "id")
		class := GetAttr(n, "class")
		if (step.Id != "" && step.Class != "" && step.Id == id && strings.Contains(class, step.Class)) || (step.Id != "" && step.Id == id) || (step.Class != "" && strings.Contains(class, step.Class)) || (step.Id == "" && step.Class == "") {
			return true
		}
	}
	return false
}
