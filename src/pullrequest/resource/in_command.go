package resource

import (
	"fmt"
	"os"
)

// InCommand is
type InCommand struct {
	github Github
}

// NewInCommand is
func NewInCommand(g Github) *InCommand {
	return &InCommand{g}
}

// Run is
func (ic *InCommand) Run(destDir string, req InRequest) (InResponse, error) {
	resp := InResponse{}

	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return resp, err
	}

	pulls, err := ic.github.ListPRs()
	if err != nil {
		return resp, err
	}

	for _, pull := range pulls {
		if pull.SHA == req.Version.Ref {
			err = ic.github.DownloadPR(destDir, pull.Number)
			if err != nil {
				return resp, err
			}

			return InResponse{
				Version: Version{Ref: req.Version.Ref},
			}, nil
		}
	}

	return resp, fmt.Errorf("version %s not found", req.Version.Ref)
}
