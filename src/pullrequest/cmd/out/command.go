package main

import (
	"fmt"
	"io"
	r "pullrequest/resource"
)

// OutCommand is
type OutCommand struct {
	github r.Github
	writer io.Writer
}

// NewOutCommand is
func NewOutCommand(g r.Github, w io.Writer) *OutCommand {
	return &OutCommand{g, w}
}

// Run is
func (oc *OutCommand) Run(sourceDir string, req r.OutRequest) (r.OutResponse, error) {
	resp := r.OutResponse{}
	params := req.OutParams

	err := oc.github.UpdatePR(sourceDir, params.Status)
	if err != nil {
		return resp, fmt.Errorf("updating pr: %+v", err)
	}

	return resp, nil
}
