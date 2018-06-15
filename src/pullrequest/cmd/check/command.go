package main

import (
	r "pullrequest/resource"

	"github.com/google/go-github/github"
)

// CheckCommand is
type CheckCommand struct {
	github r.Github
}

// NewCheckCommand is
func NewCheckCommand(g r.Github) *CheckCommand {
	return &CheckCommand{g}
}

// Run is
func (cc *CheckCommand) Run(request r.CheckRequest) ([]r.Version, error) {
	versions := []r.Version{}

	opts := &github.PullRequestListOptions{}
	pulls, err := cc.github.ListPRs(opts)
	if err != nil {
		return versions, err
	}

	if len(pulls) == 0 {
		return versions, nil
	}

	for _, pull := range pulls {
		versions = append(versions, r.Version{Ref: *pull.Head.SHA})
	}

	return versions, nil
}
