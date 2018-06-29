package resource

import (
	"fmt"
	"io/ioutil"
	"path"
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

	prNumber, err := ioutil.ReadFile(path.Join(sourceDir, params.Path, "pr_number"))
	if err != nil {
		return OutResponse{}, fmt.Errorf("reading pr_number: %+v", err)
	}

	ref, err := oc.github.UpdatePR(sourceDir, params.Status, params.Path)
	if err != nil {
		return OutResponse{}, fmt.Errorf("updating pr: %+v", err)
	}

	return OutResponse{
		Version: Version{
			Ref: ref,
			PR:  string(prNumber),
		},
	}, nil
}
