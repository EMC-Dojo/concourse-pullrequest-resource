package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	r "pullrequest/resource"
	"pullrequest/util"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <sources directory>\n", os.Args[0])
	}

	req := r.NewOutRequest()
	util.InputRequest(&req)

	sourceDir := os.Args[1]

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("constructing github client: %+v", err)
	}

	command := r.NewOutCommand(github)
	resp, err := command.Run(sourceDir, req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	util.OutputResponse(resp)
}
