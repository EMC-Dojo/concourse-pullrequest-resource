package resource

import (
	"fmt"
	"io"
)

// OutCommand is
type OutCommand struct {
	github Github
	writer io.Writer
}

// NewOutCommand is
func NewOutCommand(g Github, w io.Writer) *OutCommand {
	return &OutCommand{g, w}
}

// Run is
func (oc *OutCommand) Run(sourceDir string, req OutRequest) (OutResponse, error) {
	resp := OutResponse{}
	params := req.OutParams

	err := oc.github.UpdatePR(sourceDir, params.Status)
	if err != nil {
		return resp, fmt.Errorf("updating pr: %+v", err)
	}

	return resp, nil
}
