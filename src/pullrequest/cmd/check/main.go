package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"

	r "pullrequest/resource"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	req := r.NewCheckRequest()
	inputRequest(&req)

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("contstructing github client: %+v", err)
	}

	command := r.NewCheckCommand(github)
	resp, err := command.Run(req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	outputResponse(resp)
}

func inputRequest(request *r.CheckRequest) {
	err := json.NewDecoder(os.Stdin).Decode(request)
	if err != nil {
		log.Fatalf("reading request from stdin: %+v", err)
	}
}

func outputResponse(response []r.Version) {
	err := json.NewEncoder(os.Stdout).Encode(response)
	if err != nil {
		log.Fatalf("writing response to stdout: %+v", err)
	}
}
