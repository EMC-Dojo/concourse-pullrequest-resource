package main

import (
	"encoding/json"
	"log"
	"os"
	r "pullrequest/resource"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <sources directory>\n", os.Args[0])
	}

	req := r.NewOutRequest()
	inputRequest(&req)

	sourceDir := os.Args[1]

	github, err := r.NewGithubClient(req.Source)
	if err != nil {
		log.Fatalf("constructing github client: %+v", err)
	}

	command := NewOutCommand(github, os.Stderr)
	resp, err := command.Run(sourceDir, req)
	if err != nil {
		log.Fatalf("running command: %+v", err)
	}

	outputResponse(resp)
}

func inputRequest(req *r.OutRequest) {
	err := json.NewDecoder(os.Stdin).Decode(req)
	if err != nil {
		log.Fatalf("reading request from stdin: %+v", err)
	}
}

func outputResponse(resp r.OutResponse) {
	err := json.NewEncoder(os.Stdout).Encode(resp)
	if err != nil {
		log.Fatalf("writing response to stdout: %+v", err)
	}
}
