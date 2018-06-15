package main

import (
	"io"
	"os"
	r "pullrequest/resource"
)

// InCommand is
type InCommand struct {
	github r.Github
	writer io.Writer
}

// NewInCommand is
func NewInCommand(g r.Github, w io.Writer) *InCommand {
	return &InCommand{g, w}
}

// Run is
func (ic *InCommand) Run(destDir string, req r.InRequest) (r.InResponse, error) {
	resp := r.InResponse{}

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return resp, err
	}

	pulls, err := ic.github.ListPRs(nil)
	if err != nil {
		return resp, err
	}

	for _, pull := range pulls {
		if *pull.GetHead().SHA == req.Version.Ref {
			err = ic.github.DownloadPR(destDir, 1)
			if err != nil {
				return resp, err
			}

			resp = r.InResponse{
				Version: r.Version{Ref: req.Version.Ref},
			}
		}
	}

	return resp, nil
}
