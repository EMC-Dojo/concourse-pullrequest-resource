package resource

import (
	"fmt"
)

// OutCommand is
type OutCommand struct {
	github Github
}

// NewOutCommand is
func NewOutCommand(g Github) *OutCommand {
	return &OutCommand{g}
}

// Run is
func (oc *OutCommand) Run(sourceDir string, req OutRequest) (OutResponse, error) {
	params := req.OutParams

	ref, err := oc.github.UpdatePR(sourceDir, params.Status, params.Path)
	if err != nil {
		return OutResponse{}, fmt.Errorf("updating pr: %+v", err)
	}

	return OutResponse{Version: Version{Ref: ref}}, nil
}
