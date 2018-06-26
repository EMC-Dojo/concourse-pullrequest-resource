package resource

import (
	"io"
	"os"
)

// InCommand is
type InCommand struct {
	github Github
	writer io.Writer
}

// NewInCommand is
func NewInCommand(g Github, w io.Writer) *InCommand {
	return &InCommand{g, w}
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
		if *pull.GetHead().SHA == req.Version.Ref {
			err = ic.github.DownloadPR(destDir, pull.GetNumber())
			if err != nil {
				return resp, err
			}

			resp = InResponse{
				Version: Version{Ref: req.Version.Ref},
			}
		}
	}

	return resp, nil
}
