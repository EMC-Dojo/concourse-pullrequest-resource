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
		log.Printf("usage: %s <source directory>\n", os.Args[0])
		os.Exit(1)
	}

	req := r.NewInRequest()
	util.InputRequest(&req)

	destDir := os.Args[1]

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("constructing github client: %+v", err)
	}

	command := r.NewInCommand(github)
	resp, err := command.Run(destDir, req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	util.OutputResponse(resp)
}
